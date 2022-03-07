var form = document.getElementById("Paramform")
var element = form.elements["Frequency"]
var index = element.selectedIndex
console.log(index)
var selectedOptionVal = element.options[index].value
console.log(selectedOptionVal)


