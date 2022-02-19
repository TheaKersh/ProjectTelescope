async function displayData(dataObj) {
  console.log("Hello, world");
  var y_list = dataObj.dataSets[0].series["0:0:0:0"].observations
  var y_in = []
  console.log(y_list[0].length)
  for (let ind = 0; ind < 30; ++ind) {
    console.log(ind)
    var element = y_list[ind]
    element = element[0]
    console.log(element)
    y_in.push(element)
  }
  y_in.push(300.7)
  console.log(y_in)
  
  x_in = []

  for (let i = 1990; i < 2020; ++i){
    x_in.push(i)
  }


  var data = [
    {
      x: 
      y: y_in,
      type: 'line'
    }
  ];
  

  Plotly.newPlot('myDiv', data)}


window.onload = () => {
  fetch("/view")
    .then(response => response.json())
    .then(displayData);
}
