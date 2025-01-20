import { LocalDb } from './indexdb.js';

// class AIChatParams {
//   constructor(apiUrl, streamAPIUrl, tokenKey, resultDiv,contextType,contextValue, modelNodeId, promptId, input, temperature, topP, modelName) {
//     this.apiUrl = apiUrl;
//     this.streamAPIUrl = streamAPIUrl;
//     this.tokenKey = tokenKey;
//     this.resultDiv = resultDiv;
//     this.modelNodeId = modelNodeId;
//     this.promptId = promptId;
//     this.input = input;
//     this.temperature = temperature;
//     this.topP = topP;
//     this.modelName = modelName;
//     this.contextType = contextType;
//     this.contextValue = contextValue;
//   }
// }


// const escapeHTML = (str) =>
//   str.replace(/&/g, "&amp;")
//     .replace(/</g, "&lt;")
//     .replace(/>/g, "&gt;")
//     .replace(/"/g, "&quot;")
//     .replace(/'/g, "&#039;");

function escapeHTML(str) {
  str = str.replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
    return decodeURIComponent(originalString);
}



export async function aiInterfaceAction(apiUrl,
   streamAPIUrl,
    tokenKey,
     resultDiv,
      contextType,
       contextValue,
        modelNodeId, 
        promptId,
         input,
          temperature,
           topP,
            modelName,
            maxTokens,
            topK) {
  const isStream = getToggleStreamValue() === "true";
  const reqId = generateUUID();
  const qReqId = `q-${reqId}`;
  const rReqId = `r-${reqId}`;
  const pRReqId = `pr-${reqId}`;
  let userQElem = `<div class="mb-4 flex justify-end" id=${qReqId}>
            <div class="text-right">
              <p class="text-sm font-semibold text-blue-600">You:</p>
              <p class="bg-gray-100 p-2 rounded-lg text-gray-700">
                 <pre class="bg-gray-100 text-gray-700 p-2 rounded-lg mt-2 break-words whitespace-pre-wrap">
                 ${input}
                 </pre>
              </p>
            </div>
          </div>`
  let responseElem = `
    <div class="chat-loader" id="loader-${pRReqId}">
      Replying<span>.</span><span>.</span><span>.</span>
    </div>
    <div class="mb-4 flex justify-start hidden" id=${rReqId}>
    <div class="text-left">
        <p class="text-sm font-semibold text-green-600">Vapus AI Studio:</p>
        <p class="bg-gray-800 text-white p-2 rounded-lg text-gray-700" id=${pRReqId}>
        </p>
    </div>
    </div>`
  console.log(resultDiv)
  let opBox = document.getElementById(resultDiv)
  opBox.innerHTML += userQElem;
  opBox.innerHTML += responseElem;
  opBox.scrollTop = opBox.scrollHeight;
  const radios = Array.from(document.getElementsByName('aiInteractionMode'));

  // Find the checked radio button
  const modeSel = radios.find(radio => radio.checked);
  const mode = modeSel.value;

  const myHeaders = new Headers();
  myHeaders.append("Accept", "application/x-ndjson");
  myHeaders.append("Content-Type", "application/x-ndjson");
  const apiToken = getCookie(tokenKey);
  myHeaders.append(
    "Authorization",
    `Bearer ${apiToken}`);

  let payload = {
  }
  if (promptId) {
    payload.promptId = promptId;
  }
  if (input) {
    payload.messages = [
      {
        role: "USER",
        content: input
      }
    ];
  }
  if (temperature) {
    payload.temperature = temperature;
  }
  if (topP) {
    payload.topP = topP;
  }
  if (maxTokens) {
    console.log(maxTokens)
    payload.maxOutputTokens = maxTokens;
  } 
  if (contextType && contextValue) {
    payload.contexts = [{
      key: contextType,
      value: contextValue
    }];
  }
  if (modelName) {
    payload.model = modelName;
  }
  if (modelNodeId) {
    payload.modelNodeId = modelNodeId;
  }
  payload.mode = mode;
  if (isStream) {
    return streamServer(payload, streamAPIUrl, myHeaders, resultDiv, rReqId, pRReqId);
  } else {
    return fetchServer(payload, apiUrl, myHeaders, resultDiv, rReqId, pRReqId);
  }

}

async function streamServer(payload, apiUrl, myHeaders, resultDiv, messageDivId, pDivId) {
  const raw = JSON.stringify(payload);
  console.log(pDivId,"-------------");
  console.log(messageDivId,"============");
  const response = await fetch(
    apiUrl,
    {
      method: "POST",
      headers: myHeaders,
      body: raw,
      redirect: "follow",
    }
  );

  if (!response.body) {
    console.error("ReadableStream not supported in this environment");
    return;
  }

  const reader = response.body.getReader();
  const decoder = new TextDecoder("utf-8");
  const streamedContentDiv = document.getElementById(resultDiv);
  const messageDiv = document.getElementById(messageDivId);
  messageDiv.classList.remove('hidden');
  const pDiv = document.getElementById(pDivId);
  console.log(pDiv,"============");
  let contentpre = document.createElement('pre')
  contentpre.classList.add('text-white', 'bg-gray-900', 'p-2', 'rounded-lg', 'mt-2', "break-words", "whitespace-pre-wrap")
  pDiv.appendChild(contentpre);
  document.getElementById('loader-' + pDivId).style.display="none";

  let done = false;
  while (!done) {
    const { value, done: readerDone } = await reader.read();
    if (value) {
      let decodedValue = decoder.decode(value);
      decodedValue = decodedValue.trim();
      decodedValue = decodedValue.replace(/}{/g, "},{");
      decodedValue = decodedValue.replace(/}\n{/g, "},{");
      try {
        let strval = "[" + decodedValue + "]";
        let objVal = JSON.parse(strval);
        // let cont = escapeHTML(objVal.result.output.content);
        // let cont = marked.parse(objVal.result.output.content);
        for (let i = 0; i < objVal.length; i++) {
          objVal[i].result.choices.forEach(val => {
            contentpre.innerHTML += decodeHTMLEntities(val.messages.content);
          })
          // contentpre.innerHTML += decodeHTMLEntities(objVal[i].result.output.content);
        }
        // contentpre.textContent += escapeHTML(objVal.result.output.content);
      } catch (e) {
        console.log(e);
      }
      // Append content directly to the div
      streamedContentDiv.scrollTop = streamedContentDiv.scrollHeight; // Auto-scroll to the bottom
    }
    done = readerDone;
  }
}

function decodeHTMLEntities(text) {
  const textarea = document.createElement("textarea");
  textarea.innerHTML = text;
  const decodedText = textarea.value;
  textarea.remove(); // Optional, ensures cleanup
  return decodedText;
}
async function fetchServer(payload, apiUrl, myHeaders, resultDiv, messageDivId, pDivId) {
  const raw = JSON.stringify(payload);
  try {
    const response = await fetch(
      apiUrl,
      {
        method: "POST",
        headers: myHeaders,
        body: raw,
        redirect: "follow",
      }
    );

    // if (!response.ok) {
    //   const errorData = await response.text();
    //   errMessage = handleInappResponseError(errorData);
    //   showAlert(AlertError, "AIStudio CHat error", errMessage);
    //   return 
    // }

    const jsonResponse = await response.json();
    const streamedContentDiv = document.getElementById(resultDiv);

    const messageDiv = document.getElementById(messageDivId);
    messageDiv.classList.remove('hidden');
    const pDiv = document.getElementById(pDivId);
    let contentpre = document.createElement('pre')
    contentpre.classList.add('text-white', 'bg-gray-900', 'p-2', 'rounded-lg', 'mt-2', "break-words", "whitespace-pre-wrap")
    document.getElementById('loader-' + pDivId).style.display="none";
    // Assuming the response has the content in a similar structure
    if (
      jsonResponse.choices.length > 0
    ) {
      jsonResponse.choices.forEach(val => {
        contentpre.innerHTML += marked.parse(decodeHTMLEntities(val.messages.content));
      });
    } else {
      contentpre.textContent = "No content found in response!";
    }
    pDiv.appendChild(contentpre);
    // Append content directly to the div
    // Scroll to the bottom if needed
    streamedContentDiv.scrollTop = streamedContentDiv.scrollHeight;
  } catch (error) {
    console.error("Error fetching data:", error);
  }
}

export async function crawlUrlWithContent(urlDiv, scrap = false) {
  let url = document.getElementById(urlDiv).value;
  if (scrap) {
    try {
      // If scraping flag is enabled, adjust headers or options as needed
      const response = await fetch(url, {
        method: 'GET',
        headers: {
          'Content-Type': 'text/html',
        },
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch URL: ${response.statusText}`);
      }

      // Get the content as text
      const htmlContent = await response.text();
      console.log('Crawled Content:', htmlContent);

      const tempDiv = document.createElement('div');
      tempDiv.innerHTML = htmlContent;

      // Extract plain text from the parsed HTML
      const plainText = tempDiv.textContent || tempDiv.innerText || '';

      console.log('Crawled Plain Text Content:', plainText);
      document.getElementById('contextValue').value = plainText;
      // If the content needs further parsing (e.g., extracting JSON), handle here
      return url; // Returns the raw HTML content of the URL
    } catch (error) {
      console.error('Error while crawling the URL:', error.message);
      showAlert(AlertError, "Url - " + url, "Failed to crawl the URL. Please check the console for details.");
      return null;
    }
  } else {
    document.getElementById('contextValue').value = url;
    document.getElementById('contextType').value = "URL Reference";
    return url;
  }
}

// Placeholder functions for upload options
export function uploadFromComputer() {
  const fileInput = document.createElement('input');
  fileInput.type = 'file';
  fileInput.accept = '*/*'; // Accept any file type
  fileInput.style.display = 'none';

  fileInput.addEventListener('change', (event) => {
    const files = event.target.files;
    if (files.length > 0) {
      const file = files[0];
      const reader = new FileReader();
      // File content read as text
      reader.onload = (e) => {
        const fileContent = e.target.result;
        try {
          // Attempt to parse as JSON
          const jsonData = JSON.parse(fileContent);
          document.getElementById('contextValue').value = JSON.stringify(jsonData, null, 2);
        } catch (error) {
          // If parsing fails, treat as a plain string
          document.getElementById('contextValue').value = fileContent;
        }
      };

      // Read the file content as text
      reader.readAsText(file);
    }
  });

  // Trigger the file input
  fileInput.click();
}

const AIStudioDatabase = {
  dbName: "AIStudioDB",
  storeName: "StudioPrompts",
  pKey: "userStudioAuthId",
};

let localDbConn = new LocalDb(AIStudioDatabase.dbName, AIStudioDatabase.storeName, AIStudioDatabase.pKey);
window.addEventListener("DOMContentLoaded", () => {
  console.log("Page loaded. Initialilizing IndexDB...");
  localDbConn.openIndexDB();
  console.log("IndexDB initialized successfully!");
});


export function addContextLocally(obj, isPrompt) {
  try {
    const payload = {
      userStudioAuthId: obj.userId + "_" + obj.Organization,
      value: [{
        id: generateUUIDv4(),
        userId: obj.userId,
        Organization: obj.Organization,
        timestamp: Date.now(),
        userPrompt: isPrompt,
        content: obj.content,
      }]
    }

    localDbConn.patchData(payload, payload.userStudioAuthId);
  } catch (error) {
    console.error("Error storing data locally:", error);
  } finally {
    console.log("Data stored locally!");
  }
}

window.addContextLocally = addContextLocally;
window.crawlUrlWithContent = crawlUrlWithContent;
window.uploadFromComputer = uploadFromComputer;
window.aiInterfaceAction = aiInterfaceAction;