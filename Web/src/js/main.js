function requestAll() {
  var xmlHttpReg = new XMLHttpRequest();
  if (null == xmlHttpReg) {
    alert("New XMLHttpRequest failed!");
    return;
  }
  xmlHttpReg.open("GET", "http://127.0.0.1:6688/v1/aecm/queryall", true);
  xmlHttpReg.send(null);
  xmlHttpReg.onreadystatechange = onQueryAll;

  function onQueryAll() {
    clearTable();

    if(4 == xmlHttpReg.readyState) {
      if (200 == xmlHttpReg.status && xmlHttpReg.responseText.length > 0) {
        console.log(xmlHttpReg.responseText)
        parseJosn(xmlHttpReg.responseText);
      } else {
        alert("Server error：" + xmlHttpReg.status);
      }
    }
  }

  function parseJosn(jsonStr){
    console.log("parseJson");
    if (jsonStr.length <= 0) {
      console.log("jsonString length is 0.");
      return;
    }
    jsonObj = JSON.parse(jsonStr);
    for (var v in jsonObj) {
      console.log(v + " : " + jsonObj[v].model);
      var raw = document.createElement("tr");
      var modelCel = document.createElement("td");
      modelCel.innerHTML = jsonObj[v].model;
      raw.appendChild(modelCel);
      var brandCell = document.createElement("td");
      brandCell.innerHTML = jsonObj[v].brand;
      raw.appendChild(brandCell);
      var osVersionCell = document.createElement("td");
      osVersionCell.innerHTML = jsonObj[v].osVersion;
      raw.appendChild(osVersionCell);
      var sdkVersionCell = document.createElement("td");
      sdkVersionCell.innerHTML = jsonObj[v].sdkVersion;
      raw.appendChild(sdkVersionCell);
      var packageNameCell = document.createElement("td");
      packageNameCell.innerHTML = jsonObj[v].packageName;
      raw.appendChild(packageNameCell);
      var authorCell = document.createElement("td");
      authorCell.innerHTML = jsonObj[v].author;
      raw.appendChild(authorCell);
      var insertTimeCell = document.createElement("td");
      insertTimeCell.innerHTML = jsonObj[v].insertTime;
      raw.appendChild(insertTimeCell);
      var buttonCell = document.createElement("td");
      buttonCell.innerHTML = "<input type='button' value='删除' onclick='removeMobile(this)'>";
      raw.appendChild(buttonCell);
      document.getElementById("mobles_table").appendChild(raw);
    }
  }

  function clearTable() {
    var tb = document.getElementById('mobile_table');
    var rowNum=tb.rows.length;
    for (i=1;i<rowNum;i++) {
      tb.deleteRow(i);
      rowNum=rowNum-1;
      i=i-1;
    }
  }
}

function removeMobile(obj) {
  console.log("removeMobile!" + obj.rowIndex);
  alert(obj.cells[1].value);
}
