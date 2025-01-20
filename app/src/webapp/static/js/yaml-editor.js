let editor;

// Initialize CodeMirror
function initCodeMirror() {
    editor = CodeMirror(document.getElementById("yamlEditor"), {
        mode: "yaml",
        lineNumbers: true,
        lineWrapping: true,
        theme: "default", // You can change the theme if desired
        viewportMargin: Infinity, // Makes the editor auto-expand to fit the content
        lint: true,
    });
}

// Show loading overlay
function showLoading() {
    document.getElementById("loading-overlay").classList.remove("hidden");
}

// Hide loading overlay
function hideLoading() {
    document.getElementById("loading-overlay").classList.add("hidden");
}

// Function to handle API call and show loading spinner
async function VapusDataAct() {
    mParams = document.getElementById("modalParams").innerText;
    params = JSON.parse(mParams);
    showLoading();  // Show loading overlay
    yamlContent = getYAMLEditorVal();  // Get the YAML content from the editor
    let jsonPayload;
    try {
        jsonPayload = jsyaml.load(yamlContent);
        console.log('JSON:', jsonPayload);
    } catch (error) {
        showAlert(AlertError, "VapusData Action", "Invalid YAML content");
        hideLoading();
        return;
    }
    console.log('JSON:', jsonPayload);
    try {
        // Make the API call
        const apiToken = getCookie(params.tokenKey);
        const myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        myHeaders.append("Authorization", `Bearer ${apiToken}`);
        console.log('Token:', apiToken);

        const payload = JSON.stringify(jsonPayload);
        console.log('Raw:', myHeaders);
        console.log('Raw:', payload);
        const requestOptions = {
            method: "POST",
            headers: myHeaders,
            body: payload,
            redirect: "follow"
        };
        const response = await fetch(params.apiUrl, requestOptions);
        // Process the response
        if (!response.ok) {
            const errorData = await response.text();
            // If the response is JSON, you can parse it (optional)
            message = handleInappResponseError(errorData);
            showAlert(AlertError, "VapusData Action", message);
            // throw new Error('Action failed');
        } else {
            const result = await response.json();
            hideLoading();  // Hide the loading overlay
            showAlert(AlertInfo, "VapusData Action", "Executed action successfully");
            location.reload(true);
        }

    } catch (error) {
        console.log(error);
        showAlert(AlertError, "VapusData Action", "Execution failed due to an error");
    } finally {
        hideLoading();
    }
}

function getYAMLEditorVal() {
    return editor.getValue();
}

// Load YAML file content into CodeMirror
function loadFileContent(event) {
    const file = event.target.files[0];
    if (file) {
        const reader = new FileReader();
        reader.onload = function (e) {
            editor.setValue(e.target.result);
        };
        reader.readAsText(file);
    }
}

// Initialize CodeMirror when the modal opens
function openYAMLedModal(apiUrl, tokenKey, contentDiv) {
    console.log(apiUrl);
    document.getElementById("yamlModal").classList.remove("hidden");
    params = JSON.stringify({ apiUrl: apiUrl, tokenKey: tokenKey });
    document.getElementById("modalParams").innerText = params;
    console.log(document.getElementById("modalParams").innerText);

    if (!editor) {
        initCodeMirror();
    }
    if (contentDiv != null) {
        val = document.getElementById(contentDiv).innerText;
        if (val) {
            editor.setValue(val);
        }
    }
}

function openYAMLedModalWithText(apiUrl, tokenKey, text) {
    console.log(apiUrl);
    document.getElementById("yamlModal").classList.remove("hidden");
    params = JSON.stringify({ apiUrl: apiUrl, tokenKey: tokenKey });
    document.getElementById("modalParams").innerText = params;
    console.log(document.getElementById("modalParams").innerText);

    if (!editor) {
        initCodeMirror();
    }
    editor.setValue(text);
}

function closeYAMLedModal() {
    document.getElementById("yamlModal").classList.add("hidden");
    if (editor) {
        editor.setValue(""); // Clear the editor content
    }
}
// Adjust the textarea height based on content
function adjustYamlFieldHeight(textarea) {
    textarea.style.height = "auto";  // Reset the height
    textarea.style.height = (textarea.scrollHeight) + "px";  // Set the new height based on content
}