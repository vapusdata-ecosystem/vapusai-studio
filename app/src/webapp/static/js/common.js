const AlertSuccess = "success";
const AlertError = "error";
const AlertInfo = "info";

const currentUri = encodeURIComponent(window.location.href);

const copySvg = `
<svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    width="20"
    height="20"
    fill="none"
    stroke="currentColor"
    aria-hidden="true"
>
    <path
        d="M16 2H8a2 2 0 00-2 2v1H5a2 2 0 00-2 2v11a2 2 0 002 2h4v1a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-1V4a2 2 0 00-2-2zM8 4h8v1H8V4zm12 3v11H9v-1h7a2 2 0 002-2V7h2zM5 6h11v10H5V6z"
    />
</svg>
`;

function toggleSidebar() {
    alert("here");
    const sidebar = document.getElementById("sidebar");
    sidebar.classList.toggle("collapsed");
    alert(sidebar.classList)
}

function copyToClipboard(text) {
    // const el = document.createElement('textarea');
    // el.value = text;
    // document.body.appendChild(el);
    // el.select();
    navigator.clipboard.writeText(text);
    const toast = document.getElementById('toast');
    toast.textContent = "copied to clipboard";
    toast.classList.add('show');

    // Hide the toast after 2 seconds
    setTimeout(() => {
        toast.classList.remove('show');
    }, 1000);
    // document.body.removeChild(el);
}

function showErrorMessage(header, text) {
    const toast = document.getElementById('errorMessage');
    toast.textContent = header + ": " + text;
    toast.classList.add('show');

    // Hide the toast after 2 seconds
    setTimeout(() => {
        toast.classList.remove('show');
    }, 1000);
    // document.body.removeChild(el);
}

function showInfoMessage(header, text) {
    const toast = document.getElementById('infoMessage');
    toast.textContent = header + ": " + text;
    toast.classList.add('show');

    // Hide the toast after 2 seconds
    setTimeout(() => {
        toast.classList.remove('show');
    }, 1000);
    // document.body.removeChild(el);
}




function copyToClipboardUsingElement(el) {
    text = document.getElementById(el).innerHTML;
    navigator.clipboard.writeText(text);
}

function getRandomColor() {
    const letters = '0123456789ABCDEF';
    let color = '#';
    for (let i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
}

function downloadElementIntoYAML(id, name) {
    // Get the text content you want to download
    const text = document.getElementById(id).innerText;
    // Convert text to a Blob object
    const blob = new Blob([text], { type: "text/yaml" });
    // Create a link element for download
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `${name}.yaml`; // Set the filename for the download

    // Trigger the download
    link.click();

    // Clean up by revoking the object URL
    URL.revokeObjectURL(link.href);
    document.body.removeChild(link);
}

function dataToYAML(data, name) {
    // Convert text to a Blob object
    const blob = new Blob([data], { type: "text/yaml" });

    // Create a link element for download
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `${name}.yaml`; // Set the filename for the download

    // Trigger the download
    link.click();

    // Clean up by revoking the object URL
    URL.revokeObjectURL(link.href);
    document.body.removeChild(link);
}

function dataToJSON(data, name) {
    const jsonString = JSON.stringify(data, null, 2);
    // Convert text to a Blob object
    const blob = new Blob([jsonString], { type: "application/json" });

    // Create a link element for download
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `${name}.json`; // Set the filename for the download

    // Trigger the download
    link.click();

    // Clean up by revoking the object URL
    URL.revokeObjectURL(link.href);
    document.body.removeChild(link);
}

function dataToCSV(data, name) {
    // Convert text to a Blob object
    if (name == null) {
        name = Date.now();
    }
    const content = convertToCSV(data);
    const blob = new Blob([content], { type: "text/csv" });

    // Create a link element for download
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `${name}.csv`; // Set the filename for the download

    // Trigger the download
    link.click();

    // Clean up by revoking the object URL
    URL.revokeObjectURL(link.href);
    document.removeChild(link);
}

function convertToCSV(data) {
    const csvRows = [];
    const headers = Object.keys(data[0]);
    csvRows.push(headers.join(","));

    for (const row of data) {
        const values = headers.map(header => {
            const escaped = ("" + row[header]).replace(/"/g, '\\"');
            return `"${escaped}"`;
        });
        csvRows.push(values.join(","));
    }

    return csvRows.join("\n");
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

// Function to show the custom alert
function showAlert(type, title, message) {
    // Remove existing alert if it exists
    const existingAlert = document.getElementById("custom-alert");
    if (existingAlert) {
        existingAlert.remove();
    }

    // Create the alert container
    const alertContainer = document.createElement("div");
    alertContainer.id = "custom-alert";
    alertContainer.className = "fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50";

    // Create the alert box
    const alertBox = document.createElement("div");
    alertBox.id = "alert-box";
    alertBox.className = "bg-white w-auto p-6 rounded-lg shadow-lg border-l-4";

    // Set the border color based on the type
    switch (type) {
        case AlertSuccess:
            alertBox.classList.add("border-green-500");
            break;
        case AlertError:
            alertBox.classList.add("border-red-500");
            break;
        case AlertInfo:
        default:
            alertBox.classList.add("border-blue-500");
            break;
    }

    // Create the alert header with title
    const header = document.createElement("div");
    header.className = "flex justify-between items-center mb-4";

    const alertTitle = document.createElement("h2");
    alertTitle.id = "alert-title";
    alertTitle.className = "text-xl font-semibold";
    alertTitle.textContent = title;

    // Close button
    const closeButton = document.createElement("button");
    closeButton.className = "text-gray-500 hover:text-gray-700";
    closeButton.onclick = closeAlert;

    const closeIcon = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    closeIcon.className = "h-6 w-6";
    closeIcon.setAttribute("fill", "none");
    closeIcon.setAttribute("viewBox", "0 0 24 24");
    closeIcon.setAttribute("stroke", "currentColor");

    const closePath = document.createElementNS("http://www.w3.org/2000/svg", "path");
    closePath.setAttribute("stroke-linecap", "round");
    closePath.setAttribute("stroke-linejoin", "round");
    closePath.setAttribute("stroke-width", "2");
    closePath.setAttribute("d", "M6 18L18 6M6 6l12 12");
    closeIcon.appendChild(closePath);

    closeButton.appendChild(closeIcon);

    // Append title and close button to the header
    header.appendChild(alertTitle);
    header.appendChild(closeButton);

    // Create the message paragraph
    const alertMessage = document.createElement("p");
    alertMessage.id = "alert-message";
    alertMessage.className = "text-gray-700 mb-4";
    alertMessage.textContent = message;

    // Create action button
    const actionButton = document.createElement("button");
    actionButton.className = "px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600";
    actionButton.textContent = "OK";
    actionButton.onclick = closeAlert;

    // Append elements to alert box
    alertBox.appendChild(header);
    alertBox.appendChild(alertMessage);
    alertBox.appendChild(actionButton);

    // Append alert box to container
    alertContainer.appendChild(alertBox);

    // Append the alert container to the body
    document.body.appendChild(alertContainer);
}

// Function to close the custom alert
function closeAlert() {
    const alertBox = document.getElementById("custom-alert");
    if (alertBox) {
        alertBox.remove();
    }
}

function ShowConfirm(title, message, onConfirm) {
    // Remove existing confirm modal if it exists
    const existingConfirm = document.getElementById("custom-confirm");
    if (existingConfirm) {
        existingConfirm.remove();
    }

    // Create the confirm modal container
    const confirmContainer = document.createElement("div");
    confirmContainer.id = "custom-confirm";
    confirmContainer.className = "fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50";

    // Create the confirm box
    const confirmBox = document.createElement("div");
    confirmBox.className = "bg-white w-80 p-6 rounded-lg shadow-lg border-l-4 border-blue-500";

    // Create the confirm header with title
    const header = document.createElement("div");
    header.className = "flex justify-between items-center mb-4";

    const confirmTitle = document.createElement("h2");
    confirmTitle.className = "text-xl font-semibold";
    confirmTitle.textContent = title;

    // Close button (optional)
    const closeButton = document.createElement("button");
    closeButton.className = "text-gray-500 hover:text-gray-700";
    closeButton.onclick = () => confirmContainer.remove();

    const closeIcon = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    closeIcon.className = "h-6 w-6";
    closeIcon.setAttribute("fill", "none");
    closeIcon.setAttribute("viewBox", "0 0 24 24");
    closeIcon.setAttribute("stroke", "currentColor");

    const closePath = document.createElementNS("http://www.w3.org/2000/svg", "path");
    closePath.setAttribute("stroke-linecap", "round");
    closePath.setAttribute("stroke-linejoin", "round");
    closePath.setAttribute("stroke-width", "2");
    closePath.setAttribute("d", "M6 18L18 6M6 6l12 12");
    closeIcon.appendChild(closePath);

    closeButton.appendChild(closeIcon);

    // Append title and close button to the header
    header.appendChild(confirmTitle);
    header.appendChild(closeButton);

    // Create the message paragraph
    const confirmMessage = document.createElement("p");
    confirmMessage.className = "text-gray-700 mb-4";
    confirmMessage.textContent = message;

    // Create Yes button
    const yesButton = document.createElement("button");
    yesButton.className = "px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 mr-2";
    yesButton.textContent = "Yes";
    yesButton.onclick = () => {
        onConfirm();
        confirmContainer.remove();
    };

    // Create No button
    const noButton = document.createElement("button");
    noButton.className = "px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600";
    noButton.textContent = "No";
    noButton.onclick = () => confirmContainer.remove();

    // Append elements to confirm box
    confirmBox.appendChild(header);
    confirmBox.appendChild(confirmMessage);
    confirmBox.appendChild(yesButton);
    confirmBox.appendChild(noButton);

    // Append confirm box to container
    confirmContainer.appendChild(confirmBox);

    // Append the confirm container to the body
    document.body.appendChild(confirmContainer);
}


// Close confirm function
function closeConfirm() {
    const confirmContainer = document.getElementById("custom-confirm");
    if (confirmContainer) {
        confirmContainer.remove();
    }
}

function getRequestOptions(tokenKey, method, jsPayload) {
    const apiToken = getCookie(tokenKey);
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", `Bearer ${apiToken}`);

    const requestOptions = {
        method: method,
        headers: myHeaders,
        redirect: "follow"
    };
    if (jsPayload != null) {
        requestOptions.body = JSON.stringify(jsPayload);
    }
    return requestOptions;
}

function toggleActionDropdownMenu() {
    const dropdown = document.getElementById("actionDropdownMenu");
    dropdown.classList.toggle("hidden");
}

// Show loading overlay
function showLoading() {
    document.getElementById("loading-overlay").classList.remove("hidden");
}

// Hide loading overlay
function hideLoading() {
    document.getElementById("loading-overlay").classList.add("hidden");
}

function generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        const r = (Math.random() * 16) | 0; // Generate a random number between 0 and 15
        const v = c === 'x' ? r : (r & 0x3) | 0x8; // Generate the correct digit for 'x' or 'y'
        return v.toString(16); // Convert to hexadecimal
    });
}

function generateUUIDv4() {
    return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, c =>
        (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
    );
}

function generateColors(count, opacity = 1) {
    const colors = [];
    for (let i = 0; i < count; i++) {
        const r = Math.floor(Math.random() * 256);
        const g = Math.floor(Math.random() * 256);
        const b = Math.floor(Math.random() * 256);
        colors.push(`rgba(${r}, ${g}, ${b}, ${opacity})`);
    }
    return colors;
}

function generateColorsWithBorder(count, bgOp, borderOp) {
    const bg = [];
    const borders = [];
    for (let i = 0; i < count; i++) {
        const r = Math.floor(Math.random() * 256);
        const g = Math.floor(Math.random() * 256);
        const b = Math.floor(Math.random() * 256);
        bg.push(`rgba(${r}, ${g}, ${b}, ${bgOp})`);
        borders.push(`rgba(${r}, ${g}, ${b}, ${borderOp})`);
    }
    return {
        backgroundColors: bg,
        borderColors: borders
    };
}

function toTitleCase(str) {
    if (!str) {
        return ""
    }
    return str.toLowerCase().replace(/\b\w/g, s => s.toUpperCase());
}

function handleInappResponseError(errorData) {
    // If the response is JSON, you can parse it (optional)
    try {
        const parsedError = JSON.parse(errorData);
        return parsedError.message;
    } catch (e) {
        return errorData;
    }
}

function setHttp(link) {
    if (link.search(/^http[s]?\:\/\//) == -1) {
        link = 'http://' + link;
    }
    return link;
}


// Toggle the popup visibility
function toggleContextPopup() {
    const popup = document.getElementById('topKPopup');
    popup.classList.toggle('hidden');
}

function connectGoogleDrive() {
    showAlert(AlertInfo, "Connect to Google Drive", "Note: This feature is not yet implemented.");
}

function connectOneDrive() {
    showAlert(AlertInfo, "Connect to One Drive", "Note: This feature is not yet implemented.");
}

function toggleURLCrawlFlag() {
    const urlCrawlDiv = document.getElementById('urlCrawlDiv');
    urlCrawlDiv.classList.toggle('hidden'); // Toggle visibility
}

function loadmultiSelect(elem) {
    return new Choices(elem, {
        removeItemButton: true, // Adds a remove button to selected items
        placeholderValue: "Select options",
        searchPlaceholderValue: "Search...",
        shouldSortItems: true, // Sorts dropdown items alphabetically
        noChoicesText: 'No Agents to choose from',
        classNames: {
            placeholder: 'choices__placeholders'
        },
        maxItemText: (maxItemCount) => {
            return `Only ${maxItemCount} values can be added`;
        },
    });
}

const MultiSelectCustomClasses = {
    placeholder: 'choices__placeholders',
    containerOuter: ['choices'],
    containerInner: ['choices__inner'],
    input: ['choices__input'],
    inputCloned: ['choices__input--cloned'],
    list: ['choices__list'],
    listItems: ['choices__list--multiple'],
    listSingle: ['choices__list--single'],
    listDropdown: ['choices__list--dropdown'],
    item: ['choices__item'],
    itemSelectable: ['choices__item--selectable'],
    itemDisabled: ['choices__item--disabled'],
    itemChoice: ['choices__item--choice'],
    description: ['choices__description'],
    group: ['choices__group'],
    groupHeading: ['choices__heading'],
    button: ['choices__button'],
    activeState: ['is-active'],
    focusState: ['is-focused'],
    openState: ['is-open'],
    disabledState: ['is-disabled'],
    highlightedState: ['is-highlighted'],
    selectedState: ['is-selected'],
    flippedState: ['is-flipped'],
    loadingState: ['is-loading'],
    notice: ['choices__notice'],
    addChoice: ['choices__item--selectable', 'add-choice'],
    noResults: ['has-no-results'],
    noChoices: ['has-no-choices'],
};

function toggleSubmenu(submenuId) {
    const submenu = document.getElementById(submenuId);
    const arrow = document.getElementById(submenuId + '-Arrow');
    if (submenu.classList.contains('hidden')) {
        submenu.classList.remove('hidden'); // Show the submenu
        arrow.classList.add('rotate-180'); // Rotate the arrow
    } else {
        submenu.classList.add('hidden'); // Hide the submenu
        arrow.classList.remove('rotate-180'); // Reset arrow rotation
    }
}

function fetchAttributeFromDiv(divId, attribute) {
    const div = document.getElementById(divId);
    if (div) {
        return div.getAttribute(attribute);
    }
    return null;
}


// Function to add a new box dynamically
function addstreamDataBox(title, content, canvas, error) {
    // try {
    //     spinner = document.getElementById('spinner');
    //     spinner.remove();
    // } catch (error) {
    //     console.error(error);
    // }
    const box = document.createElement('div');
    box.classList.add('box', "flex", "h-auto", "b-4/6", "items-center", "opacity-70", "p-1", "break-words", "border-l-8", "text-gray-700");
    // if (error) {
    //     box.classList.add('bg-red-700', 'border-red-900');
    //     box.classList.remove("bg-green-700", "border-green-900");
    // }
    let fContent = "";
    try {
        var json = JSON.parse(content);
        var contentVals = Object.values(json)
        fContent = contentVals.join(", ");
    } catch (error) {
        fContent = content;
    }
    if (error) {
        box.innerHTML = `
            <span>
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="red" width="36px" height="36px">
                    <path d="M12 10.586L6.293 4.879l-1.414 1.414L10.586 12l-5.707 5.707 1.414 1.414L12 13.414l5.707 5.707 1.414-1.414L13.414 12l5.707-5.707-1.414-1.414L12 10.586z"/>
                </svg>
            </span>
            <span class="p-1 font-bold text-sm text-red-900">
                ${title} 
            </span>
            <span class="p-1 text-xs font-semibold text-red-900  break-words">
                ${fContent}
            </span>`;
    } else {
        if (title === "SUCCESSFULL") {
            box.innerHTML = `
            <span>
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="green" width="36px" height="36px">
                    <path d="M9 16.2l-3.5-3.5-1.4 1.4L9 19 20 8l-1.4-1.4z"></path>
                </svg>
            </span>
            <span class="p-1 text-sm font-bold text-green-900">
                ${title} 
            </span>
            <span class="p-1 text-xs font-semibold text-green-900  break-words">
                ${fContent}
            </span>`;
        } else {
            box.innerHTML = `
                <span>
                    ${stepSvg}
                </span>
                <span class="p-1 text-sm font-bold text-gray-700">
                    ${title} 
                </span>
                <span class="p-1 text-xs font-semibold text-gray-700  break-words">
                    ${fContent}
                </span>`;
        }
    }
    // box.innerHTML = `<p class="break-wrods"><strong>${title} : </strong> ${content} </p>`;
    box.classList.add('highlighted');

    // Append the new box to the container
    canvas.appendChild(box);
    // spinner = document.createElement('div');
    // spinner.id = 'spinner';
    // spinner.classList.add('spinner-stream');
    // canvas.appendChild(spinner);
    canvas.scrollTop = canvas.scrollHeight;

    // After a short delay, remove the highlight to create the faded effect
    setTimeout(() => {
        box.classList.remove('highlighted');
        box.classList.add('opacity-70');
    }, 1000); // Highlight lasts for 1 second
}

function addstreamDataErrBox(cont, canvas) {
    let obj = {};
    try {
        obj = JSON.parse(cont);
    } catch (error) {
        obj = {
            key: "Error",
            value: cont
        }
    }
    const box = document.createElement('div');
    box.classList.add('min-w-0', "max-w-full");
    content = `
        <div class="mt-2 w-full p-1 shadow-md rounded-lg bg-white">
            <div class="text-primary-black flex min-h-16 justify-start">
              <div class="w-1 self-stretch border px-1.5 bg-red-600 border-red-600">
              </div>
              <div class="flex grow items-left justify-between overflow-hidden p-2 bg-white">
                <div class="flex w-full gap-4">
                    <svg viewBox="0 0 24 24" class="h-8 w-8" fill="none" stroke="currentColor">
                        <path fill="currentColor" d="M11 9h2V7h-2m1 13c-4.41 0-8-3.59-8-8s3.59-8 8-8s8 3.59 8 8s-3.59 8-8 8m0-18A10 10 0 0 0 2 12a10 10 0 0 0 10 10a10 10 0 0 0 10-10A10 10 0 0 0 12 2m-1 15h2v-6h-2z">
                        </path>
                    </svg>
                  <div class="flex w-full flex-col text-red-800">
                      <p class="text-sm font-bold items-left">
                          ${obj.key}
                      </p>
                      <p class="text-xs font-medium">
                          ${obj.value}
                      </p>
                  </div>
                </div>
              <div class="ml-2">
              </div>
            </div>
          </div>
        </div>
    `
    box.innerHTML = content;
    canvas.appendChild(box);
    // spinner = document.createElement('div');
    // spinner.id = 'spinner';
    // spinner.classList.add('spinner-stream');
    // canvas.appendChild(spinner);
    canvas.scrollTop = canvas.scrollHeight;
}

function addstreamStateBox(cont, chatId, canvas) {

    hideStreamDiv(chatId);
    const box = document.createElement('div');
    box.classList.add('chat-loader');
    box.id = `loader-${chatId}`;
    box.innerHTML = `${cont}<span>.</span><span>.</span><span>.</span>`
    canvas.appendChild(box);
    canvas.scrollTop = canvas.scrollHeight;
}

function addstreamDataSuccessBox(cont, canvas) {
    let obj = {};
    try {
        obj = JSON.parse(cont);
    } catch (error) {
        console.log(error);
        obj = {
            key: "Status",
            value: cont
        }
    }
    const box = document.createElement('div');
    box.classList.add('min-w-0', "max-w-full");
    content = `
        <div class="mt-2 w-full p-1 shadow-md rounded-lg bg-white">
            <div class="text-primary-black flex min-h-16 justify-start">
              <div class="w-1 self-stretch border px-1.5 bg-green-800 border-green-800">
              </div>
              <div class="flex grow items-left justify-between overflow-hidden p-1 bg-white">
                <div class="flex w-full gap-4">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="h-6 w-6" fill="green" stroke="currentColor">
                        <path d="M9 16.2l-3.5-3.5-1.4 1.4L9 19 20 8l-1.4-1.4z"></path>
                    </svg>
                  <div class="flex w-full flex-col text-green-800">
                      <p class="text-sm font-bold items-left">
                          ${obj.key}
                      </p>
                      <p class="text-xs font-medium">
                          ${obj.value}
                      </p>
                  </div>
                </div>
              <div class="ml-2">
              </div>
            </div>
          </div>
        </div>
    `
    box.innerHTML = content;
    canvas.appendChild(box);
    // spinner = document.createElement('div');
    // spinner.id = 'spinner';
    // spinner.classList.add('spinner-stream');
    // canvas.appendChild(spinner);
    canvas.scrollTop = canvas.scrollHeight;
}

function addDataSetHead(canvas, title, text) {
    const box = document.createElement('div');
    box.innerHTML = `
    <div class="flex justify-start mb-1 justify-between">
    <button class="flex items-center px-4 py-2 bg-gray-600 text-gray-100 rounded-lg hover:bg-pink-900" 
    title="Copy to clipboard" onclick="copyToClipboard('${text}')">
    ${title}: ${copySvg}
  </button>
  </div>`;
    canvas.appendChild(box);
    canvas.scrollTop = canvas.scrollHeight;
}

function addTable(tableContainerId, tableId) {
    try {
        let tableContainer = document.getElementById(tableContainerId);
        let responseTableSet = document.createElement('table');
        responseTableSet.classList.add("bg-white", "rounded-lg", "shadow-md", "border-2", "border-gray-500");
        tableContainer.appendChild(responseTableSet);
        responseTableSet.id = tableId;
        return responseTableSet;
    } catch (error) {
        console.error(error);
        return null;
    }
}

function addTableBody(table) {
    try {
        let responseTableBody = document.createElement('tbody');
        table.appendChild(responseTableBody);
        return responseTableBody;
    } catch (error) {
        console.error(error);
        return null;
    }
}

function addTableHeader(table, columns) {
    try {
        const headerElem = document.createElement('thead');
        let headerRow = document.createElement('tr');
        columns.forEach(column => {
            let headerCell = document.createElement('th');
            headerCell.innerHTML = column;
            // headerCell.classList.add('px-2', 'py-2', 'text-left', 'text-sm', 'font-medium', 'text-gray-700', 'border-b');
            headerCell.classList.add("bg-gray-50", "text-sm", "text-left", "text-gray-600", "px-2", "py-2", "border-b");
            headerRow.appendChild(headerCell);
        });
        headerElem.appendChild(headerRow);
        table.appendChild(headerElem);
    } catch (error) {
        console.error(error);
    }
}

function addTableFooter(table, reqId) {
    try {
        let tFooter = document.createElement('tfoot');
        tFooter.classList.add("bg-white");
        tFooter.innerHTML = `
        <tr>
      <td colspan="6" class="px-2 py-2">
        <div class="flex items-center justify-between">
          <!-- Left Side: Pagination Controls -->
          <div class="flex items-center space-x-2">
            <button id="${reqId}-prevPage" class="px-2 py-2 bg-black text-sm text-white rounded-lg">Previous</button>
            <button id="${reqId}-nextPage" class="px-2 py-2 bg-black text-sm text-white rounded-lg">Next</button>
            <div class="text-sm text-gray-800">
              Showing <span id="${reqId}-currentPageSize">0</span> of <span id="${reqId}-totalElements">0</span> entries
            </div>
          </div>
          <!-- Right Side: Pagination Info -->
        </div>
      </td>
    </tr>
        `;
        table.appendChild(tFooter);
    } catch (error) {
        console.error(error);
    }
}

function addTableRow(tableBody, columns, datalist) {
    try {
        datalist.forEach(data => {
            let row = document.createElement('tr');
            console.log(data);
            columns.forEach(column => {
                let cell = document.createElement('td');
                try {
                    cell.innerHTML = JSON.parse(data[column])
                } catch (error) {
                    cell.innerHTML = data[column];
                }
                cell.classList.add('px-1', 'py-1', 'text-xs', 'text-gray-600', 'border-b', 'border-l', 'border-r');
                row.appendChild(cell);
            });
            tableBody.appendChild(row);
        });
        // let row = document.createElement('tr');
        // columns.forEach(column => {
        //     let cell = document.createElement('td');
        //     try {
        //         cell.innerHTML = JSON.parse(data[column])
        //     } catch (error) {
        //         cell.innerHTML = data[column];
        //     }
        //     cell.classList.add('px-2', 'py-2', 'text-sm', 'text-gray-600', 'border-b');
        //     row.appendChild(cell);
        // });
        // table.appendChild(row);
    } catch (error) {
        console.error(error);
    }
}

const stepSvg = `    <svg xmlns="http://www.w3.org/2000/svg" width="25" height="25" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <!-- Outer Chevron -->
      <polyline points="9 18 15 12 9 6" 
                stroke="url(#gradient1)" 
                stroke-width="2" 
                stroke-linecap="round" 
                stroke-linejoin="round" 
                fill="none"/>
      
      <!-- Inner Chevron for Depth -->
      <polyline points="9 17 14 12 9 7" 
                stroke="url(#gradient2)" 
                stroke-width="2" 
                stroke-linecap="round" 
                stroke-linejoin="round" 
                fill="none"/>
      
      <!-- Decorative Dots -->
      <circle cx="5" cy="12" r="2" fill="url(#gradient3)" />
      <circle cx="19" cy="12" r="2" fill="url(#gradient3)" />
      
      <!-- Gradient Definitions -->
      <defs>
        <linearGradient id="gradient1" x1="9" y1="18" x2="15" y2="6" gradientUnits="userSpaceOnUse">
          <stop offset="0%" stop-color="black"/>
          <stop offset="100%" stop-color="black"/>
        </linearGradient>
        <linearGradient id="gradient2" x1="9" y1="17" x2="14" y2="7" gradientUnits="userSpaceOnUse">
          <stop offset="0%" stop-color="black"/>
          <stop offset="100%" stop-color="black" stop-opacity="0.5"/>
        </linearGradient>
        <radialGradient id="gradient3" cx="12" cy="12" r="10" fx="12" fy="12">
          <stop offset="0%" stop-color="black"/>
          <stop offset="100%" stop-color="black"/>
        </radialGradient>
      </defs>
    </svg>`;

function uploadFromComputer() {
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.accept = '.csv,application/json,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'; // Accept CSV, JSON, and XLSX files
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

function strToUniArray(str) {
    const encoder = new TextEncoder(); // Defaults to UTF-8
    return encoder.encode(str);
}

function uint8ArrayToBase64(uint8Array) {
    let binary = '';
    uint8Array.forEach((byte) => {
        binary += String.fromCharCode(byte);
    });
    return btoa(binary);
}

function showTab(tabId) {
    const tabs = document.querySelectorAll('.tab-content');
    tabs.forEach(tab => tab.classList.add('hidden'));

    document.getElementById(tabId).classList.remove('hidden');

    const buttons = document.querySelectorAll('.tab-button');
    buttons.forEach(button => button.classList.remove('active-tab'));

    document.querySelector(`[onclick="showTab('${tabId}')"]`).classList.add('active-tab');
}

function getLoaderTextDiv(text, reqId) {
    return `<div class="chat-loader" id="loader-${reqId}">
      ${text}<span>.</span><span>.</span><span>.</span>
    </div>`;
}

function hideStreamDiv(reqId) {
    try {
        console.log("removing loader", `loader-${reqId}`);
        el = document.getElementById(`loader-${reqId}`);
        el.remove();
    } catch (error) {
        console.error(error);
    }

}

function addDatasetDownloadButton(tableId, canvas) {
    const box = document.createElement('div');
    box.innerHTML = `
    <button onclick="exportResultSet('${tableId}')"
    class="flex items-center px-4 py-2 bg-gray-600 text-gray-100 rounded-lg hover:bg-pink-900">
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" class="mr-2" xmlns="http://www.w3.org/2000/svg">
      <path
        d="M5 20h14v2H5v-2zm7-2c-.28 0-.53-.11-.71-.29L8 13.41l1.41-1.41L11 14.17V4h2v10.17l1.59-1.59L16 13.41l-3.29 3.29c-.18.18-.43.29-.71.29z"
        fill="#FFFFFF" />
    </svg>
    Download Report
  </button>
    `
    canvas.appendChild(box);
    canvas.scrollTop = canvas.scrollHeight;
}

function renderFiles(response, canvas) {
    console.log(response);
    if (response.format.toLowerCase() === "png") {
        renderstreamImage(response, canvas);
    } else if (response.format.toLowerCase() === "html") {
        renderhtmlPageInObj(response, canvas);
    }
}

function renderhtmlPageInObj(response, canvas) {
    const box = document.createElement("div");
    const fileUrl = `data:text/html;base64,${response.data}`;
    const object = document.createElement("object");
    object.type = "text/html";
    object.data = fileUrl;
    object.classList.add("w-2/3", "h-96");
    box.classList.add("mt-2", "p-1", "shadow-md", "rounded-lg", "bg-white");
    box.appendChild(object);
    canvas.appendChild(box);

    console.log(fileUrl);
    const nBox = document.createElement("div");
    nBox.classList.add("mt-4");
    const downloadLink = document.createElement("a");
    downloadLink.classList.add("items-center", "px-4", "py-2", "bg-gray-600", "text-gray-100", "rounded-lg", "hover:bg-pink-900");
    downloadLink.href = fileUrl;
    downloadLink.textContent = `Download`;
    downloadLink.download = response.name;
    nBox.appendChild(downloadLink);
    canvas.appendChild(nBox);
}

function renderstreamImage(response, canvas) {
    const box = document.createElement("div");
    const imageDataUrl = `data:image/${response.format.toLowerCase()};base64,${response.data}`;
    const img = document.createElement("img");
    img.src = imageDataUrl;
    img.alt = response.name;
    img.classList.add("w-2/3", "h-96");
    box.classList.add("mt-2", "p-1", "shadow-md", "rounded-lg", "bg-white");
    box.appendChild(img);
    canvas.appendChild(box);

    const nBox = document.createElement("div");
    nBox.classList.add("mt-4");
    const downloadLink = document.createElement("a");
    downloadLink.classList.add("items-center", "px-4", "py-2", "bg-gray-600", "text-gray-100", "rounded-lg", "hover:bg-pink-900");
    downloadLink.href = imageDataUrl;
    downloadLink.textContent = `Download`;
    downloadLink.download = response.name;
    nBox.appendChild(downloadLink);
    canvas.appendChild(nBox);
}


function isJSONorJSONArray(str) {
    try {
        const parsed = JSON.parse(str);

        if (typeof parsed === "object" && parsed !== null) {
            if (Array.isArray(parsed)) {
                const allObjects = parsed.every(item => typeof item === "object" && item !== null);
                return parsed, allObjects;
            } else {
                return parsed, false;
            }
        }
        return undefined, false;
    } catch (e) {
        // If parsing fails, it's not valid JSON
        return undefined, false;
    }
}


function buildAndDownloadFile(fileName, format, base64Data) {
    // Determine the MIME type based on the format
    format = format.toLowerCase();
    // let name = file.name.replace("."+format, "");
    const mimeTypes = {
        "png": "image/png",
        "csv": "text/csv",
        "json": "application/json",
        "yaml": "application/x-yaml",
    };

    // Check if the format is supported
    if (!mimeTypes[format]) {
        throw new Error(`Unsupported format: ${format}`);
    }

    // Decode the base64 data to binary
    const binaryData = atob(base64Data);
    const byteNumbers = new Uint8Array(binaryData.length);
    for (let i = 0; i < binaryData.length; i++) {
        byteNumbers[i] = binaryData.charCodeAt(i);
    }

    // Create a Blob with the decoded data and the correct MIME type
    const blob = new Blob([byteNumbers], { type: mimeTypes[format] });

    // Create a URL for the Blob
    const url = URL.createObjectURL(blob);

    // Create a temporary anchor element to trigger the download
    const a = document.createElement("a");
    a.href = url;
    a.download = `${fileName}`;
    document.body.appendChild(a);

    // Trigger the download
    a.click();

    // Cleanup: Remove the anchor element and revoke the Blob URL
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}

function getFileExtension(fileName) {
    // Find the last dot in the file name
    const lastDotIndex = fileName.lastIndexOf('.');

    // If there's no dot or it's the first character, return an empty string
    if (lastDotIndex === -1 || lastDotIndex === 0) {
        return '';
    }

    // Return the substring after the last dot
    return fileName.slice(lastDotIndex + 1);
}

async function fileUploader(apiUrl, accessTokenKey, resource, resourceId) {
    return new Promise((resolve, reject) => {
        const fileInput = document.createElement("input");
        fileInput.type = "file";
        fileInput.accept = ".csv,application/json,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"; // Accept CSV, JSON, and XLSX files
        fileInput.style.display = "none";

        fileInput.addEventListener("change", async (event) => {
            const files = event.target.files;
            if (files.length > 0) {
                const file = files[0];
                const reader = new FileReader();

                // Read the file content
                reader.onload = async (e) => {
                    const fileContent = e.target.result;
                    try {
                        // Process file data
                        let data = strToUniArray(fileContent);
                        data = uint8ArrayToBase64(data);
                        const fileExt = getFileExtension(file.name);

                        const payload = {
                            resourceId: resourceId,
                            resource: resource,
                            objects: [
                                {
                                    name: file.name,
                                    data: data,
                                    format: fileExt.toUpperCase(),
                                },
                            ],
                        };

                        const myHeaders = new Headers();
                        myHeaders.append("Accept", "application/x-ndjson");
                        myHeaders.append("Content-Type", "application/x-ndjson");
                        const apiToken = getCookie(accessTokenKey);
                        myHeaders.append("Authorization", `Bearer ${apiToken}`);

                        const params = JSON.stringify(payload);

                        // Send the file upload request
                        const response = await fetch(apiUrl, {
                            method: "POST",
                            headers: myHeaders,
                            redirect: "follow",
                            body: params,
                        });

                        if (!response.ok) {
                            showAlert(AlertError, "Upload Failed", "File upload failed, please try after some time");
                            reject("File upload failed");
                            return;
                        }

                        const resultBody = await response.json();
                        let fileResults = [];
                        if (resultBody.output.length > 0) {
                            resultBody.output.forEach((res) => {
                                console.log(res);
                                fileResults.push({
                                    name: res.object.name,
                                    format: res.object.format,
                                    responsePath: res.responsePath,
                                    fid: res.fid
                                });
                            });
                        }
                        resolve(fileResults); // Resolve the promise with the result
                    } catch (error) {
                        console.error("Error uploading file:", error);
                        showAlert(AlertError, "Upload Failed", "File upload failed, please try after some time");
                        reject([]); // Reject the promise in case of errors
                    }
                };

                // Read the file as text
                reader.readAsText(file);
            } else {
                reject("No file selected"); // Reject if no file is selected
            }
        });

        // Trigger the file input
        fileInput.click();
    });
}

function updateBrowserUrl(key, value) {
    let currentUrl = window.location.href;

      // Check if the URL already contains query parameters
      let newUrl = new URL(currentUrl);

      // Add a new query parameter (or update if exists)
      newUrl.searchParams.set(key, value);

      // Update the URL in the browser without reloading
      window.history.pushState({}, '', newUrl);
}