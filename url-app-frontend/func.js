const button = document.getElementById('post-btn');

button.addEventListener('click', async _ => {
  var textUrl = document.getElementById("text-line").value;

  console.log(textUrl);
  var xhr = new XMLHttpRequest();
  xhr.open("POST", "http://localhost:8080/post", true);
  xhr.setRequestHeader('Content-Type', 'x-www-form-urlencoded');//x-www-form-urlencoded
  xhr.send(JSON.stringify({
      "url": textUrl
  }));
});