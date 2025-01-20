import { LocalDb } from './indexdb.js';

const DatafabricDatabase = {
  dbName: "DataFabricDB",
  storeName: "VapusQueries",
  pKey: "queryId",
  batchSize: 200,
  waitTime: 100
};

const NeedAgentRun = false;
let currentChatId = '';
let currChatObj;
let selectedMessageId = '';
let attachDatasets = [];
let localDbConn = new LocalDb(DatafabricDatabase.dbName, DatafabricDatabase.storeName, DatafabricDatabase.pKey);
window.addEventListener("DOMContentLoaded", () => {
  console.log("Page loaded. Initialilizing IndexDB...");
  localDbConn.openIndexDB();
  console.log("IndexDB initialized successfully!");
});

export function toggleInputType() {
  const inputType = document.getElementById('inputFormat').value;
  const textInput = document.getElementById('textInput');

  if (inputType === 'query') {
    textInput.classList.remove('hidden');
  } else if (inputType === 'api') {
    textInput.classList.add('hidden');
  }
}

async function uploadDataset(api, tokenKey, resource, fileCanvasId,chatId) {
  if (chatId !== "") {
    currentChatId = chatId;
    await fabricFileUploader(api, tokenKey, resource, chatId, fileCanvasId);
  } else {
    if (currChatObj) {
      if (currChatObj !== null && currChatObj !== undefined) {
        currentChatId = currChatObj.fabricChatId;
        await fabricFileUploader(api, tokenKey, resource, currentChatId, fileCanvasId);
      } else {
        showAlert(AlertError, "Error", "Please start a chat to upload a dataset");
      }
    } else {
      showAlert(AlertError, "Error", "Please start a chat to upload a dataset");
    }
  }
  
}

async function createNewChat(manageUrl, tokenKey) {
  const myHeaders = new Headers();
  myHeaders.append("Accept", "application/x-ndjson");
  myHeaders.append("Content-Type", "application/x-ndjson");
  const apiToken = getCookie(tokenKey);
  myHeaders.append(
    "Authorization",
    `Bearer ${apiToken}`);
  const payload = JSON.stringify({
    action: "CREATE"
  });
  try {
    const response = await fetch(
      manageUrl,
      {
        method: "POST",
        headers: myHeaders,
        body: payload,
        redirect: "follow",
      }
    );
    if (!response.ok) {
      showErrorMessage("Error", "Error while creating new chat, please try again");
      return;
    } else {
      console.log("Chat created successfully");
      const result = await response.json();
      console.log(result);
      if (result.output !== undefined && result.output.length === 1) {
        return result.output[0];
      } else {
        showErrorMessage("Error", "Error while creating new chat, please try again");
      }
    }
  } catch (error) {
    console.error("Error creating new chat:", error);
    showErrorMessage("Error", "Error while creating new chat, please try again");
    return null;
  }
}
export async function dataFabricAction(tokenKey, apiUrl, manageUrl, username, chId) {
  const input = document.getElementById('input').value;
  document.getElementById('input').disabled = true;
  const dpId = document.getElementById('dataProduct').value;
  if (input === "") {
    showErrorMessage("Data Product Missing", "Please enter a query");
    return;
  }
  const genSug = document.getElementById('fabric-suggetion-generic');
  genSug.classList.add("hidden");
  console.log(currChatObj, "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++");
  if (chId !== "") {
    currentChatId = chId;
  } else {
    if (currChatObj !== null && currChatObj !== undefined) {
      currentChatId = currChatObj.fabricChatId;
    } else {
      console.log("Creating new chat");
      currentChatId = '';
      currChatObj = await createNewChat(manageUrl, tokenKey);
      console.log(currChatObj);
      if (currChatObj === null || currChatObj === undefined) {
        showErrorMessage("Error", "Error while creating new chat, please try again");
        return;
      } else {
        currentChatId = currChatObj.fabricChatId;
        updateBrowserUrl("fabricChatId", currentChatId);
      }
    }
  }
  var queryParams = {
    fabricChatId: currentChatId,
    input: input,
  };
  if (dpId === "") {
    queryParams.dataproducts = [];
  } else {
    queryParams.dataproducts = [dpId];
  }
  if (selectedMessageId === "") {
    queryParams.messageId = [];
  } else {
    queryParams.messageId = [selectedMessageId];
  }
  if (attachDatasets.length > 0) {
    queryParams.fileData = [];
    attachDatasets.forEach(dataset => {
      queryParams.fileData.push({
        name: dataset
      });
    });
  }
  selectedMessageId = '';
  const myHeaders = new Headers();
  myHeaders.append("Accept", "application/x-ndjson");
  myHeaders.append("Content-Type", "application/x-ndjson");
  const apiToken = getCookie(tokenKey);
  myHeaders.append(
    "Authorization",
    `Bearer ${apiToken}`);
  const payload = JSON.stringify(queryParams);
  // requestOptions.mode = "no-cors";
  // if (streamChat === 'true') {
  return submitQueryForStream(apiUrl, myHeaders, payload, username);
  // } else {
  //   fetch(url, requestOptions)
  //     .then((response) => response.text())
  //     .then((result) => displayResponseTable(result))
  //     .catch((error) => handleNoResponseOrErr(error));
  // }
}
async function submitQueryForStream(url, myHeaders, payload, username) {
  try {
    let fabricCanvas = document.getElementById("dataServerOutput");
    fabricCanvas.classList.remove("hidden");
    const messageCanvasElem = document.createElement("div");
    messageCanvasElem.classList.add("my-2");
    const tempMessid = generateUUID();
    messageCanvasElem.id = tempMessid;
    fabricCanvas.appendChild(messageCanvasElem);
    const messageCanvas = document.getElementById(tempMessid);
    addFabricUserMessage(username, document.getElementById('input').value, messageCanvas);
    document.getElementById('input').value = "";
    const response = await fetch(
      url,
      {
        method: "POST",
        headers: myHeaders,
        body: payload,
        redirect: "follow",
      }
    );

    if (!response.body) {
      console.error("ReadableStream not supported in this environment");
      return;
    }
    const reader = response.body.getReader();
    const decoder = new TextDecoder("utf-8");
    let done = false;

    let dataFields = [];
    let messCahe;
    let counter = 0;
    let responseId = "";
    let dataTable, dataTableBody;

    while (!done) {
      messageCanvas.scrollTop = messageCanvas.scrollHeight;
      fabricCanvas.scrollTop = fabricCanvas.scrollHeight;
      const { value, done: readerDone } = await reader.read();
      if (value) {
        let decodedValue = decoder.decode(value);
        decodedValue = decodedValue.trim();
        try {
          decodedValue = decodedValue.replace(/}{/g, "},{");
          decodedValue = decodedValue.replace(/}\n{/g, "},{");
          let strval;
          let objVal;
          try {
            strval = "[" + decodedValue + "]";
            objVal = JSON.parse(strval);
            if (objVal.length < 0) {
              continue;
            }
            messCahe = null;
          } catch (error) {
            if (messCahe !== null) {
              messCahe = messCahe + decodedValue;
            } else {
              messCahe = decodedValue;
            }
            try {
              strval = "[" + messCahe + "]";
              objVal = JSON.parse(strval);
            } catch (error) {
              continue;
            }

            console.error("Error parsing JSON: WIll concartenate with prvious messages", error);
          }

          objVal.forEach(responseJson => {
            hideStreamDiv(responseId);
            if (responseJson !== null && responseJson.error !== undefined) {
              addstreamDataErrBox(responseJson.error.message, messageCanvas);
              return;
            }
            if (responseJson === null || responseJson.result === null || responseJson.result.output === null) {
              addstreamDataErrBox("Error while querying data product either there is no data for this query or there is some internal server error, please try again or contact the data product owner",
                messageCanvas);
              return;
            }
            switch (responseJson.result.output.event) {
              case "END":
                if (responseJson.result.output.data.final.reason === "SUCCESSFULL") {
                  addstreamDataSuccessBox(responseJson.result.output.data.content, messageCanvas);
                } else if (responseJson.result.output.data.final.reason === "DONE") {
                  break;
                } else {
                  addstreamDataErrBox(JSON.stringify({
                    key: responseJson.result.output.data.final.reason,
                    value: responseJson.result.output.data.final.metadata
                  }), messageCanvas);
                }
                break;
              case "ABORTED":
                addstreamDataErrBox(JSON.stringify({
                  key: responseJson.result.output.data.final.reason,
                  value: responseJson.result.output.data.final.metadata
                }), messageCanvas);
                break;
              case "START":
                setChatId(responseJson.result.fabricChatId);
                break;
              case "DATASET_END":
                const datadiv = document.createElement("div");
                datadiv.classList.add("overflow-x-auto", "w-auto", "border-gray-800");
                datadiv.appendChild(dataTable);
                messageCanvas.appendChild(datadiv);
                addTableFooter(dataTable, responseId);
                addstreamDataBox("Download", "You can select the file format for the report and then click on the download button below", messageCanvas, false);
                addDatasetDownloadButton(responseId, messageCanvas);
                break;
              case "DATASET_START":
                // addDataSetHead(messageCanvas, "DataSet Id: ", responseId);
                dataTable = addTable("dataServerOutput", responseId);
                dataTableBody = addTableBody(dataTable);
                dataFields = JSON.parse(responseJson.result.output.data.content);
                addTableHeader(dataTable, dataFields);
                break;
              case "STATE":
                addstreamStateBox(responseJson.result.output.data.content, currentChatId, messageCanvas);
                break;
              case "FILE_DATA":
                renderFiles(responseJson.result.output.files, messageCanvas);
                break;
              case "RESPONSE_ID":
                responseId = responseJson.result.responseId;
                messageCanvas.id = responseId;
                break;
              default:
                if (responseJson.result.output.data.contentType === "JSON") {
                  let content = JSON.parse(responseJson.result.output.data.content);
                  addstreamDataBox(content.key, content.value, messageCanvas, false);
                } else if (responseJson.result.output.data.contentType === "DATASET") {
                  counter = handleDataset(responseJson.result.output.data.dataset, dataTableBody, dataFields, responseId, counter);
                  counter++;
                } else {
                  addstreamDataBox("Info", responseJson.result.output.data.content, messageCanvas, false);
                }
                break;
            }
          });

        } catch (error) {
          console.error(error);
          addstreamDataErrBox("Error while processing your request, please try again", messageCanvas);
        } finally {
          hideStreamDiv(responseId);
        }
      }
      done = readerDone;
    }
    try {
      hideStreamDiv(responseId);
      hideStreamDiv(currentChatId);


      const divSperator = document.createElement("div");
      divSperator.classList.add("flex", "mt-2", "cursor-pointer", "text-xs", "text-gray-400", "hover:text-gray-600");
      divSperator.innerHTML = "Click here to select this input again in you chat";
      const innerDiv = document.createElement("div");
      innerDiv.innerHTML = copySvg;
      innerDiv.addEventListener("click", function () {
        selectedMessageId = messageCanvas.id;
        const toast = document.getElementById('toast');
        toast.textContent = "message selected";
        toast.classList.add('show');
      });
      const seperator = document.createElement("hr");
      seperator.classList.add("border-gray-400", "border-2", "my-2");
      divSperator.appendChild(innerDiv);
      messageCanvas.appendChild(divSperator);
      messageCanvas.appendChild(seperator);
      fabricCanvas.appendChild(messageCanvas);
      setTimeout(() => {
        addTablePagination({
          tableId: dataTable.id,
          rowsPerPage: 10,
          prevPageBtn: document.getElementById(responseId + "-prevPage"),
          nextPageBtn: document.getElementById(responseId + "-nextPage"),
          currentPageSizeSpan: document.getElementById(responseId + "-currentPageSize"),
          totalElementsSpan: document.getElementById(responseId + "-totalElements"),
        });
        // addDataSetInIndexDb(responseId, dbBatch);
      }, 2 * DatafabricDatabase.waitTime);

    } catch (error) {
      console.error(error);
    }
  } catch (error) {
    console.error("Error:", error);
  } finally {
    document.getElementById('input').disabled = false;
    // attachDatasets = [];
  }
}

function setChatId(op) {
  if (op !== "") {
    currentChatId = op;
    return;
  }
}


function handleDataset(dataset, dataTableBody, dataFields, responseId, counter) {
  let dbBatch = [];
  let rowbatch = [];
  if (counter < DatafabricDatabase.batchSize) {
    dbBatch.push(dataset);
    rowbatch.push(dataset);
  }
  setTimeout(() => {
    addTableRow(dataTableBody, dataFields, rowbatch);
    addDataSetInIndexDb(responseId, dbBatch);
  }, DatafabricDatabase.waitTime);
  return 0;
}

function addDataSetInIndexDb(iden, dataSet) {
  try {
    const payload = {
      queryId: iden,
      value: dataSet
    }

    localDbConn.patchData(payload, iden);
  } catch (error) {
    console.error("Error storing dataset locally:", error);
  }
}

export async function exportResultSet(responseId) {
  const dataset = await localDbConn.retrieveData(responseId);
  const format = document.getElementById('exportFormat').value;
  // let parsedData = JSON.parse(dataset);
  // console.log(parsedData);
  if (format === 'JSON') {
    dataToJSON(dataset.value, '');
  } else if (format === 'CSV') {
    dataToCSV(dataset.value, '');
  } else {
    console.error('Invalid export format');
  }
}

export function showPrompts() {
  try {
    const dp = document.getElementById("dataProduct").value;
    if (dp === "") {
      document.getElementById("dp-prompts").classList.add("hidden");
      showErrorMessage("Invalid param", "Please select a data product");
      return;
    } else {
      document.getElementById("prompts-" + dp).classList.remove("hidden");
    }
  } catch (error) {
    showErrorMessage("Invalid param", "Please select a data product");
  }
}

export function closePrompts(id) {
  document.getElementById(id).classList.add("hidden");
}


function selectPromptIntoInput(prompt, inputDiv, promptId) {
  try {
    const input = document.getElementById(inputDiv);
    input.value = prompt;
    var event = new KeyboardEvent('keydown', {
      key: 'Enter',
      code: 'Enter',
      keyCode: 13,
      which: 13,
      bubbles: true,
    });

    input.dispatchEvent(event);
    closePrompts(promptId)

  } catch (error) {
    console.error("Error selecting prompt into input:", error);
  }
}

function loadFabricChat(username, accessTokenKey, downloadUrl, canvasId, currFabChatObject) {
  try {
    console.log("Loading fabric chat", currChatObj);
    if (currFabChatObject !== null) {
      const canvas = document.getElementById(canvasId);
      canvas.classList.remove("hidden");
      console.log(currFabChatObject, currFabChatObject.messages);
      if (currFabChatObject.messages !== undefined && currFabChatObject.messages.length > 0) {
        currFabChatObject.messages.forEach(message => {
          console.log(message);
          console.log(message.input);
          const messageCanvas = document.createElement("div");
          messageCanvas.classList.add("border-gray-400", "my-2");
          messageCanvas.id = message.messageId;
          if (message.input === undefined) {
            if (message.error !== undefined && message.error !== "") {
              addstreamDataErrBox(message.error, messageCanvas);
            } else {
              return;
            }
          }
          addFabricUserMessage(username, message.input.content, messageCanvas);
          if (message.error !== undefined && message.error !== "") {
            addstreamDataErrBox(message.error, messageCanvas);
          }
          if (message.dataproducts !== undefined && message.dataproducts.length > 0) {
            addstreamDataBox("Selected data products", message.dataproducts.join(","), messageCanvas, false);
          }
          if (message.datafields !== undefined && message.datafields.length > 0) {
            addstreamDataBox("Selected data fields", message.datafields.join(","), messageCanvas, false);
          }
          if (message.datasetLength !== undefined && message.datasetLength !== "") {
            addstreamDataBox("Dataset Length", "total data points in result: " + message.datasetLength, messageCanvas, false);
          }
          if (message.output !== undefined && message.output.length > 0) {
            message.output.forEach(output => {
              addstreamDataBox("Output", output.content, messageCanvas, false);
              if (output.file !== undefined && output.file.length > 0) {
                // renderstreamImage(output.file, canvas);
                output.file.forEach(file => {
                  if (file !== undefined && file !== "" && file.name !== undefined && file.name !== "") {
                    addFabricFileDownloadButton(messageCanvas, file, "click here to download the file", accessTokenKey, downloadUrl);
                  }
                });
              }
            });
          }
          const divSperator = document.createElement("div");
          divSperator.classList.add("flex", "mt-2", "cursor-pointer", "text-xs", "text-gray-400", "hover:text-gray-600");
          divSperator.innerHTML = "Click here to select this input again in you chat";
          const innerDiv = document.createElement("div");
          innerDiv.innerHTML = copySvg;
          innerDiv.addEventListener("click", function () {
            selectedMessageId = message.messageId;
            const toast = document.getElementById('toast');
            toast.textContent = "message selected";
            toast.classList.add('show');
          });
          const seperator = document.createElement("hr");
          seperator.classList.add("border-gray-400", "border-2", "my-2");
          divSperator.appendChild(innerDiv);
          messageCanvas.appendChild(divSperator);
          messageCanvas.appendChild(seperator);
          canvas.appendChild(messageCanvas);
        });
        if (currFabChatObject.ended) {
          document.getElementById("userInputGround").classList.add("hidden");
          document.getElementById("endedChatMessage").classList.remove("hidden");
        }
      }
      canvas.scrollTop = canvas.scrollHeight;
    }
  } catch (error) {
    console.error("Error loading fabric chat:", error);
  }
}

function addFabricFileDownloadButton(canvas, file, text, accessTokenKey, downloadUrl) {
  const nBox = document.createElement("div");
  nBox.classList.add("m-4");
  const downloadLink = document.createElement("a");
  downloadLink.classList.add("items-center", "text-xs", "px-4", "py-2", "bg-gray-600", "text-gray-100", "rounded-lg", "hover:bg-pink-900");
  downloadLink.href = "#";
  downloadLink.textContent = text;
  downloadLink.onclick = () => downloadFabricFile(file.name, accessTokenKey, downloadUrl);
  nBox.appendChild(downloadLink);
  canvas.appendChild(nBox);
  canvas.scrollTop = canvas.scrollHeight;
}

async function downloadFabricFile(fileName, accessTokenKey, downloadUrl) {
  const myHeaders = new Headers();
  myHeaders.append("Accept", "application/x-ndjson");
  myHeaders.append("Content-Type", "application/x-ndjson");
  const apiToken = getCookie(accessTokenKey);
  myHeaders.append(
    "Authorization",
    `Bearer ${apiToken}`);
  const payload = JSON.stringify({ fileName: fileName });
  downloadUrl = downloadUrl + "?fileName=" + fileName;
  try {
    const response = await fetch(downloadUrl, {
      method: "GET",
      headers: myHeaders,
      redirect: "follow",
    });
    if (!response.ok) {
      showAlert(AlertError, "Download Failed", "File download failed, please try after some time");
    }

    const result = await response.json();
    result.output.forEach(file => {
      if (file.name !== undefined && file.name !== "") {
        buildAndDownloadFile(file.name, file.format, file.data);
      }

    });
  } catch (error) {
    console.error("Error downloading file:", error);
    showAlert(AlertError, "Download Failed", "File download failed, please try after some time");
  }
}



async function fabricFileUploader(apiUrl, tokenKey, resource, resourceId, fileCanvasId) {
  const fileCanvas = document.getElementById(fileCanvasId);
  const result = await fileUploader(apiUrl, tokenKey, resource, resourceId);
  if (result.length > 0) {
    result.forEach(file => {
      if (file.name !== undefined && file.name !== "") {
        addFabricFile(file, fileCanvas);
      }
    });
  }
}

function addFabricFile(file, fileCanvas) {
  const liItem = document.createElement("li");
  liItem.classList.add("flex", "p-1", "items-center", "justify-between");
  const filenamediv = document.createElement("div");
  filenamediv.textContent = file.name + "attach";
  filenamediv.id = file.fid + "-name";
  filenamediv.classList.add("p-1", "m-1", "bg-white", "border-1", "text-orange-800", "border-black",
    "text-xs", "shadow-sm", "rounded-md", "overflow-hidden", "whitespace-nowrap");
  const attachDiv = document.createElement("div");
  attachDiv.innerHTML = `<svg viewBox="0 0 24 24" width="1.2em" height="1.2em" class="size-3"><path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6z"></path></svg>`;
  attachDiv.classList.add("p-1", "m-1", "bg-orange-800", "border-1", "text-white", "border-black", "hover:bg-pink-900", "hover:text-white", "cursor-pointer");

  const fileInfo = document.createElement("input");
  fileInfo.type = "hidden";
  fileInfo.value = file.responsePath;
  fileInfo.id = file.fid;
  attachDiv.onclick = function () {
    attachDataset(this, file.fid);

  };

  const detachdiv = document.createElement("div");
  detachdiv.innerHTML = `<svg viewBox="0 0 24 24" width="1.2em" height="1.2em" class="size-3">
  <path fill="currentColor" d="M5 11h14v2H5z"></path>
</svg>`;
  detachdiv.classList.add("hidden", "p-1", "m-1", "text-orange-800", "border-1", "bg-white", "border-black", "hover:bg-pink-900", "hover:text-white", "cursor-pointer");
  detachdiv.onclick = function () {
    dettachDataset(this, file.fid);
  };
  detachdiv.id = file.fid + "-detach";
  attachDiv.id = file.fid + "-attach";

  liItem.appendChild(fileInfo);
  liItem.appendChild(filenamediv);
  liItem.appendChild(attachDiv);
  liItem.appendChild(detachdiv);
  fileCanvas.appendChild(liItem);
}

function attachDataset(el, filesetId) {
  const fileDiv = document.getElementById(filesetId + "-name");
  if (fileDiv === null) {
    return;
  }
  fileDiv.classList.add("fileuploader-list-selected");
  const filePath = document.getElementById(filesetId);
  if (filePath === null) {
    return;
  }
  const existIndex = attachDatasets.indexOf(filePath.value);
  if (existIndex < 0) {
    attachDatasets.push(filePath.value);
    showInfoMessage("File Attached", "File attached successfully");
  }
  el.classList.add("hidden");
  const detachdiv = document.getElementById(filesetId + "-detach");
  detachdiv.classList.remove("hidden");
}

function dettachDataset(el, filesetId) {
  const fileDiv = document.getElementById(filesetId + "-name");
  if (fileDiv === null) {
    return;
  }
  fileDiv.classList.remove("fileuploader-list-selected");
  const filePath = document.getElementById(filesetId);
  if (filePath === null) {
    return;
  }
  const existIndex = attachDatasets.indexOf(filePath.value);
  if (existIndex > -1) {
    attachDatasets.pop(filePath.value);
    showInfoMessage("File Attached", "File attached successfully");
  }
  el.classList.add("hidden");
  const attachDiv = document.getElementById(filesetId + "-attach");
  attachDiv.classList.remove("hidden");
}


function addFabricUserMessage(user, content, canvas) {
  let seperator = document.createElement("div");
  seperator.style.margin = "4px 0";
  seperator.classList.add("flex", "justify-end");
  seperator.innerHTML = `<div class="bg-gray-200 font-semibold break-words px-4 py-2 rounded-lg mt-3 mb-2 w-full">
                          <p class="text-xs text-blue-700 mt-1">${user}:</p>
                          <p class="text-sm text-gray-700-500"> ${content}</p>
                      </div>`;
  canvas.appendChild(seperator);
}

function clearAttachDataset(id) {
  attachDatasets = [];
}

window.toggleInputType = toggleInputType;
window.dataFabricAction = dataFabricAction;
window.exportResultSet = exportResultSet;
window.showPrompts = showPrompts;
window.closePrompts = closePrompts;
window.selectPromptIntoInput = selectPromptIntoInput;
window.loadFabricChat = loadFabricChat;
window.fabricFileUploader = fabricFileUploader;
window.currChatObj = currChatObj;
window.uploadDataset = uploadDataset;
