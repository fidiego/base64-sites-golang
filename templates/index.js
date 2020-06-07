const generateResults = function(response) {
  // remove results, if exists
  var results = document.getElementById('results');
  if (results) {
    results.remove();
  }
  var resultLinks = document.getElementById('results-links');
  if (resultLinks) {
    results.remove();
  }
  // generate result nodes
  var section = document.getElementById('try-it');
  var node = document.createElement('PRE');
  node.id = 'results';
  node.textContent = JSON.stringify(response, null, '\t');
  section.appendChild(node);
  var links = document.createElement('div');
  links.id = 'results-link';
  links.innerHTML =
    '<p>Copy <code>base64_string</code> (or <a href="' +
    response.base64_string +
    '">this</a> link) and paste it into your URL bar.</p>' +
    '<p>Or use <a href="/render?content=' +
    encodeURIComponent(response.base64_string) +
    '">this link</a> on a mobile device.</p>';
  // append results
  section.appendChild(links);
};

const onSubmit = function() {
  var ta = document.getElementById('content-textarea');
  var html = ta.value;
  var payload = {content: html};
  fetch('/api', {
    method: 'POST',
    headers: {'content-type': 'application/json'},
    body: JSON.stringify(payload),
  })
    .then(function(response) {
      response.json().then(function(resp) {
        generateResults(resp);
      });
    })
    .catch(function(err) {
      console.error(err);
    });
};

(function() {
  console.info('Document Ready: registering Event Listener');
  var form = document.getElementById('content-form');
  form.addEventListener('submit', function(e) {
    e.preventDefault();
    onSubmit();
  });
})();
