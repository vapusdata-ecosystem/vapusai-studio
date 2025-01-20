function EnterInput(event) {
    if (event.key === "Enter") {
        submitInput();
    }
}
var AgentStepsContext = [];
var StepsEnum = [];
window.addEventListener("DOMContentLoaded", () => {
    try{
    vals = document.getElementById("stepsEnums").innerText;
    StepsEnum = JSON.parse(vals);
    }catch(error){
        console.error(error);
    }
});


// Toggle the popup visibility
function openAgentContextPopup() {
    const aiAgent = document.getElementById("aiAgent").value;
    const stepsData = document.getElementById(aiAgent + "-steps");
    if (!stepsData) {
        showAlert(AlertError, "Agent Steps", "Please select an agent first");
        return;
    }
    const steps = JSON.parse(stepsData.innerText);
    container = document.getElementById("agentStepsContainer");
    stepsElem = document.createElement("select");
    stepsElem.classList.add("w-full", "px-3", "py-2", "bg-blue-950", "text-white", "rounded-lg", "focus:outline-none", "focus:ring", "focus:ring-blue-200");
    stepsElem.setAttribute("id", "agent-steps-dropdown");
    steps.forEach(step => {
        option = document.createElement("option");
        option.value = StepsEnum[step.id];
        option.text = StepsEnum[step.id];
        stepsElem.appendChild(option);
    });
    container.appendChild(stepsElem);
    const popup = document.getElementById('topKPopup');
    popup.classList.toggle('hidden');
}

function closeAgentContextPopup() {
    content = document.getElementById("agent-steps-dropdown");
    let data = document.getElementById("contextValue").value;
    data = strToUniArray(data);
    data = uint8ArrayToBase64(data);

    const step = { stepId: content.value, data: data };
    AgentStepsContext.push(step);
    const popup = document.getElementById('topKPopup');
    document.getElementById("agentStepsContainer").innerHTML = "";
    document.getElementById("contextValue").value = "";
    popup.classList.toggle('hidden');
}

function showModels(id) {
    const models = document.getElementById(`models-` + id);
    modelElems = document.getElementsByClassName("modelList");
    for (let i = 0; i < modelElems.length; i++) {
        modelElems[i].classList.add("hidden");
    }
    models.classList.remove("hidden");
}
function showDatastores(id) {
    const models = document.getElementById(`datastores-` + id);
    modelElems = document.getElementsByClassName("datastores");
    for (let i = 0; i < modelElems.length; i++) {
        modelElems[i].classList.add("hidden");
    }
    models.classList.remove("hidden");
}
// function showSteps(id) {
//     const agentId = document.getElementById(id).value;
//     if (!agentId) {
//         showAlert(AlertError, "Agent ID", "Please select an agent first");
//         return;
//     }
//     const steps = document.getElementById(`` + agentId + `-steps`).innerText;
//     const stepsObj = JSON.parse(steps);
//     console.log(stepsObj);
//     stepsEnum = {};
//     try {
//         se = document.getElementById("stepsEnum").innerText;
//         stepsEnum = JSON.parse(se);
//     } catch (error) {
//         console.log(error);
//     }
//     console.log(stepsEnum);
//     for (const [key, value] of Object.entries(stepsObj)) {
//         let step = value.id.toString();
//         if (stepsEnum[step]) {
//             appendAgentSteps(stepsEnum[step], "agent-steps", value.prompt);
//         } else {
//             appendAgentSteps(step, "agent-steps", value.prompt)
//         }
//     }
// }
// function appendAgentSteps(inputId, parentDivId, placeholder) {
//     // Create the outer div container for the input field and label
//     const containerDiv = document.createElement('div');
//     containerDiv.classList.add('mb-4');
//     let labelText = inputId;
//     console.log(labelText);
//     console.log(inputId);
//     console.log(placeholder);
//     try {
//         labelText = labelText.replace(`ST_`, '');
//         labelText = labelText.replace(/_/g, ' ');
//         labelText = toTitleCase(labelText);
//     } catch (error) {
//         console.log(error);
//     }
//     // Create the label element
//     const label = document.createElement('label');
//     label.setAttribute('for', inputId);
//     label.classList.add('block', 'text-gray-700', 'font-medium', 'mb-1', 'agent-steps');
//     label.textContent = labelText;

//     // Create the input element
//     const input = document.createElement('input');
//     input.setAttribute('type', 'text');
//     input.setAttribute('id', inputId);
//     input.setAttribute('placeholder', placeholder);
//     input.classList.add('w-full', 'px-3', 'py-2', 'border', 'rounded-lg', 'focus:outline-none', 'focus:ring', 'focus:ring-blue-200');

//     // Append label and input to the container
//     containerDiv.appendChild(label);
//     containerDiv.appendChild(input);

//     // Find the parent element by its ID and append the new container
//     const parentDiv = document.getElementById(parentDivId);
//     if (parentDiv) {
//         parentDiv.appendChild(containerDiv);
//     } else {
//         console.error('Parent div not found');
//     }
// }

async function agentAction(api, streamAPI,
    accessTokenKey, resultDivId,
    agentId, textInput, dataproduct, stepsContext) {
    try {
        const agentCanvas = document.getElementById(resultDivId);

        console.log(stepsContext, "====================");
        let seperator = document.createElement("div");
        seperator.style.margin = "4px 0";
        seperator.classList.add("flex", "justify-start");
        seperator.innerHTML = `<div class="bg-gray-200 font-semibold break-words px-4 py-2 rounded-lg mb-2 max-w-full">
                            <p class="text-xs text-blue-700 mt-1">Agent Input:</p>
                            <p class="text-sm text-gray-700-500">`+ textInput + `</p>
                        </div>`;
        agentCanvas.appendChild(seperator);

        // const runHeader = document.createElement('div');
        // runHeader.classList.add('flex', 'w-full', 'items-center', 'bg-blue-950', 'my-2', 'rounded-lg');
        // runHeader.innerHTML = `
        //         <span class="mx-4 text-gray-100 text-lg font-bold">Agent Thread: </span>
        // `;
        // agentCanvas.appendChild(runHeader);
        if (dataproduct!=="") {
        textInput = textInput + "\nData Product id: " + dataproduct;
        } else {
            textInput = textInput+"\n IMPORTANT: Dataset is already is rerady in process so data product and is not required";
        }
        queryParams = {
            agentId: agentId,
            input: textInput,
            steps: stepsContext,
        };
        console.log(stepsContext);
        // if (textInput === "") {
        //     queryParams.steps = [];
        //     agentSteps = document.getElementsByClassName("agent-steps");
        //     for (let i = 0; i < agentSteps.length; i++) {
        //         console.log(agentSteps[i]);
        //         let step = agentSteps[i].id;
        //         let value = document.getElementById(step).value;
        //         queryParams.steps.push({ stepId: step, value: value });
        //         console.log({ stepId: step, value: value });
        //     }
        // }
        payload = {
            chain: [queryParams],
        }
        requestOptions = getRequestOptions(accessTokenKey, "POST", payload);
        const response = await fetch(
            streamAPI,
            requestOptions
        );

        if (!response.body) {
            console.error("ReadableStream not supported in this environment");
            return;
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder("utf-8");
        let done = false;
        AgentStepsContext = [];
        while (!done) {
            const { value, done: readerDone } = await reader.read();
            if (value) {
                let decodedValue = decoder.decode(value);
                decodedValue = decodedValue.trim();
                try {
                    decodedValue = decodedValue.replace(/}{/g, "},{");
                    decodedValue = decodedValue.replace(/}\n{/g, "},{");
                    strval = "[" + decodedValue + "]";
                    objVal = JSON.parse(strval);
                    console.log(objVal);
                    if (objVal.length < 0) {
                        continue;
                    }
                    objVal.forEach(responseJson => {
                        if (responseJson.result.output.event === "END") {
                            if (responseJson.result.output.data.final.reason === "SUCCESSFULL") {
                                addstreamDataSuccessBox(JSON.stringify({
                                    key:responseJson.result.output.data.final.reason,
                                    value: "Yay!..Agent executed tasks successfully"
                                }), agentCanvas);
                            } else {
                                addstreamDataErrBox(JSON.stringify({
                                    key: responseJson.result.output.data.final.reason,
                                    value: responseJson.result.output.data.final.metadata
                                }), agentCanvas);
                            }
                        } else {
                            if (responseJson.result.output.data.contentType === "JSON") {
                                content = JSON.parse(responseJson.result.output.data.content);
                                addstreamDataBox(content.key, content.value, agentCanvas, false);
                            } else {
                                addstreamDataBox("Steps", responseJson.result.output.data.content, agentCanvas, false);
                            }
                        }
                    });
                } catch (error) {
                    console.error(error);
                    addstreamDataErrBox("Error: Internal Server error ,agent run failed.", agentCanvas, true);
                    try {
                        spinner = document.getElementById('spinner');
                        spinner.remove();
                    } catch (error) {
                        console.error(error);
                    }
                    // showAlert(AlertError, "Agent failed", "An error occurred while agent was processing your request");
                }
            }
            done = readerDone;
        }
        try {
            spinner = document.getElementById('spinner');
            spinner.remove();
        } catch (error) {
            console.error(error);
        }
    } catch (error) {
        console.error(error);
        try {
            spinner = document.getElementById('spinner');
            spinner.remove();
        } catch (error) {
            console.error(error);
        }
    }
}

