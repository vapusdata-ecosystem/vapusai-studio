import { LocalDb } from './indexdb.js';

const DatafabricDatabase = {
  dbName: "DataFabricDB",
  storeName: "VapusQueries",
  pKey: "queryId",
  batchSize: 200,
  waitTime: 100
};

const NeedAgentRun = false;
let currentRequestId = '';

let localDbConn = new LocalDb(DatafabricDatabase.dbName, DatafabricDatabase.storeName, DatafabricDatabase.pKey);
window.addEventListener("DOMContentLoaded", () => {
  console.log("Page loaded. Initialilizing IndexDB...");
  localDbConn.openIndexDB();
  console.log("IndexDB initialized successfully!");
});

function clearSelectedAgent() {
  const buttons = document.querySelectorAll('.button-stack');
  buttons.forEach(button => {
    button.classList.remove('button-stack-selected');
  });
  document.getElementById('vapusAgents').value = "";
}

export function toggleInputType() {
  const inputType = document.getElementById('inputFormat').value;
  const textInput = document.getElementById('textInput');
  const jsonEditorContainer = document.getElementById('jsonEditorContainer');

  if (inputType === 'query') {
    textInput.classList.remove('hidden');
    jsonEditorContainer.classList.add('hidden');
  } else if (inputType === 'api') {
    textInput.classList.add('hidden');
    jsonEditorContainer.classList.remove('hidden');
  }
}

export async function invokeAgent(tokenKey, agentServerUrl, agentId) {
  clearSelectedAgent();
  try {
    const dataset = await localDbConn.retrieveData(currentRequestId);
    console.log("Dataset:", dataset);
    if (!dataset) {
      console.error("No dataset found for requestId:", currentRequestId);
      showErrorMessage("No Datat", "Please run a query first to get the data and then select an agent to run on the data");
      return;
    }
    let data = JSON.stringify(dataset.value);
    data = strToUniArray(data);
    data = uint8ArrayToBase64(data);
    const textInput = document.getElementById('textInput').value;
    const dataproduct = "";
    const stepsContext = [
      {
        stepId: "AGENTST_DATASET",
        data: data,
      }
    ];
    agentAction(agentServerUrl, agentServerUrl,
      tokenKey, "dataServerOutput",
      agentId, textInput, dataproduct, stepsContext)

  } catch (error) {
    console.error("Error invoking agent:", error);
  } finally {
    clearSelectedAgent();
    document.getElementById('textInput').value = '';
  }
}

function handleResponseErr(errMessage) {
  const responseTableContainer = document.getElementById('responseTableContainer');
  const responseTableBody = document.getElementById('responseTableBody');
  // Clear previous table rows
  responseTableBody.innerHTML = '';
  const row = document.createElement('tr');

  row.innerHTML = `
        <td class="px-4 py-2 text-gray-700 border-b">${errMessage}</td>
    `;
  responseTableBody.appendChild(row);
  // Show the table container
  responseTableContainer.classList.remove('hidden');
}

export function dataFabricAction(tokenKey, apiUrl, agentServerUrl) {
  let agentSelected = document.getElementById('vapusAgents').value;
  console.log("Agent Selected:", agentSelected);
  if (agentSelected !== "") {
    if (currentRequestId === "") {
      showErrorMessage("Missing Input", "Please enter your message first, based on data you can select the agent");
      clearSelectedAgent();
      return;
    }
    invokeAgent(tokenKey, agentServerUrl, agentSelected);
    return;
  }
  const inputType = document.getElementById('inputFormat').value;
  const textInput = document.getElementById('textInput').value;
  const dataTables = document.getElementById('dataTablesInput').value.split(',').map(item => item.trim());
  const limit = parseInt(document.getElementById('limitInput').value, 10) || 0;
  const orderField = document.getElementById('orderFieldInput').value;
  const orderBy = document.getElementById('orderByInput').value;
  const filters = JSON.parse(document.getElementById('filtersInput').value || '[]'); // expecting JSON array of filter objects
  const columns = document.getElementById('columnsInput').value.split(',');
  const dpId = document.getElementById('dataProduct').value;
  clearSelectedAgent();
  let datastore = '';
  // Build the request payload
  var queryParams = {
    dataproducts: dpId,
    databases: datastore,
    apiQueryParam: {
      dataTables: dataTables,
      limit: limit,
      columns: columns,
      orderBy: orderBy,
      orderField: orderField,
      filters: filters
    }
  };
  if (inputType === 'text') {
    queryParams.textQuery = textInput;
  } else if (inputType === 'query') {
    // Add JSON input to the request payload
    queryParams.rawQuery = textInput;
  }

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
  return submitQueryForStream(apiUrl, myHeaders, payload);
  // } else {
  //   fetch(url, requestOptions)
  //     .then((response) => response.text())
  //     .then((result) => displayResponseTable(result))
  //     .catch((error) => handleNoResponseOrErr(error));
  // }
}
async function submitQueryForStream(url, myHeaders, payload) {
  try {
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
    const requestId = generateUUID();
    currentRequestId = requestId;
    console.log("Request ID:==================", currentRequestId);
    const reader = response.body.getReader();
    const decoder = new TextDecoder("utf-8");
    let done = false;
    let fabricCanvas = document.getElementById("dataServerOutput");
    let dataFields = [];
    let seperator = document.createElement("div");
    seperator.style.margin = "4px 0";
    seperator.classList.add("flex", "justify-start");
    seperator.innerHTML = `<div class="bg-gray-200 font-semibold break-words px-4 py-2 rounded-lg mb-2 max-w-xs">
                            <p class="text-xs text-blue-700 mt-1">You:</p>
                            <p class="text-sm text-gray-700-500">`+ document.getElementById('textInput').value + `</p>
                        </div>`;
    fabricCanvas.appendChild(seperator);
    document.getElementById('textInput').value = '';
    const dataTable = addTable("dataServerOutput", requestId);
    let downloadableReport = false;
    let messCahe;
    let counter = 0;
    while (!done) {
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
            if (responseJson !== null && responseJson.error !== undefined) {
              addBox("Error", responseJson.error.message, fabricCanvas, true);
              return;
            }
            if (responseJson === null || responseJson.result === null || responseJson.result.output === null) {
              handleResponseErr("Error while querying data product either there is no data for this query or there is some internal server error, please try again or contact the data product owner");
            }
            if (responseJson.result.output.event === "END") {
              if (responseJson.result.output.data.final.reason === "SUCCESSFULL") {
                addBox(responseJson.result.output.data.final.reason, responseJson.result.output.data.final.metadata, fabricCanvas, false);
              } else {
                addBox(responseJson.result.output.data.final.reason,
                  responseJson.result.output.data.final.metadata, fabricCanvas, true);
              }
            } else if(responseJson.result.output.event === "DATA_END") {
              console.log("Data End event received");
              const datadiv = document.createElement("div");
              datadiv.classList.add("overflow-x-auto", "w-auto","border-2","border-gray-800");
              datadiv.appendChild(dataTable);
              fabricCanvas.appendChild(datadiv);
              addTableFooter(dataTable, requestId);
              console.log("Table footer added");
            } else if(responseJson.result.output.event === "DATA_START") {
              dataFields = JSON.parse(responseJson.result.output.data.content);
              addTableHeader(dataTable, dataFields);
            } else {
              if (responseJson.result.output.data.contentType === "JSON") {
                let content = JSON.parse(responseJson.result.output.data.content);
                addBox(content.key, content.value, fabricCanvas, false);
              } else if (responseJson.result.output.data.contentType === "DATASET") {
                downloadableReport = true;
                // if (!headerSet) {
                //   dataFields = Object.keys(responseJson.result.output.data.dataset);
                //   addTableHeader(dataTable, dataFields);
                //   headerSet = true;
                // }
                counter = handleDataset(responseJson.result.output.data.dataset, dataTable, dataFields, requestId, counter);
                counter++;
              } else {
                addBox("Info", responseJson.result.output.data.content, fabricCanvas, false);
              }
            }
          });
        } catch (error) {
          console.error(error);
          addBox("Error", "Error while rendering large data", fabricCanvas, true);
          try {
            spinner = document.getElementById('spinner');
            spinner.remove();
          } catch (error) {
            console.error(error);
          }
        }
      }
      if (readerDone && downloadableReport) {
        try {
          spinner = document.getElementById('spinner');
          spinner.remove();

        } catch (error) {
          console.error(error);
        }
        // let downloadLink = document.createElement("button");
        // downloadLink.id = requestId;
        // downloadLink.classList.add("flex","items-center","px-4","py-2","bg-black","text-gray-100","rounded-lg","hover:bg-gray-800");
        // downloadLink.onclick = exportResultSet(requestId);
        addBox("Download", "You can select the file format for the report and then click on the download button below", fabricCanvas, false);
        let dwElem = `
        <button onclick="exportResultSet('${requestId}')"
            class="flex items-center px-4 py-2 bg-gray-600 text-gray-100 rounded-lg hover:bg-green-600">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" class="mr-2" xmlns="http://www.w3.org/2000/svg">
              <path
                d="M5 20h14v2H5v-2zm7-2c-.28 0-.53-.11-.71-.29L8 13.41l1.41-1.41L11 14.17V4h2v10.17l1.59-1.59L16 13.41l-3.29 3.29c-.18.18-.43.29-.71.29z"
                fill="#FFFFFF" />
            </svg>
            Download Report
          </button>
            `;
        fabricCanvas.innerHTML += dwElem;
        addBox("Agent Support", "You can select agents to run on this data by clicking on the agent icon on the right sidebar and provide the input data.", fabricCanvas, false);
      }
      done = readerDone;
    }
    try {
      spinner = document.getElementById('spinner');
      spinner.remove();
      console.log("Spinner removed");
      console.log(dataTable);
      console.log(requestId);
      addTablePagination({
        tableId: requestId,
        rowsPerPage: 10,
        prevPageBtn: document.getElementById(requestId + "-prevPage"),
        nextPageBtn: document.getElementById(requestId + "-nextPage"),
        currentPageSizeSpan: document.getElementById(requestId + "-currentPageSize"),
        totalElementsSpan: document.getElementById(requestId + "-totalElements"),
      });
    } catch (error) {
      console.error(error);
    }
  } catch (error) {
    try {
      spinner = document.getElementById('spinner');
      spinner.remove();
    } catch (error) {
      console.error(error);
    }
    console.error("Error:", error);
  }
}


function handleDataset(dataset, dataTable, dataFields, requestId, counter) {
  let dbBatch = [];
  let rowbatch = [];
  if (counter < DatafabricDatabase.batchSize) {
    dbBatch.push(dataset);
    rowbatch.push(dataset);
  }
  setTimeout(() => {
    addTableRow(dataTable, dataFields, rowbatch);
    addDataSetInIndexDb(requestId, dbBatch);
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

export async function exportResultSet(requestId) {
  const dataset = await localDbConn.retrieveData(requestId);
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
  dp = document.getElementById("dataProduct").value;
  if (dp === "") {
    document.getElementById("dp-prompts").classList.add("hidden");
    return;
  } else {
    document.getElementById("prompts-" + dp).classList.remove("hidden");
  }
}

export function closePrompts(id) {
  document.getElementById(id).classList.add("hidden");
}

window.toggleInputType = toggleInputType;
window.dataFabricAction = dataFabricAction;
window.exportResultSet = exportResultSet;
window.showPrompts = showPrompts;
window.closePrompts = closePrompts;



