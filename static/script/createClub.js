function submitForm(event) {
  event.preventDefault();
  console.log("Submitting form!");

  const url = "/api/clubCreate/";
  const data = {
    name: document.getElementById("name").value,
    description: document.getElementById("description").value
  };

  fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data)
  })
    .then(response => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.json();
    })
    .then(data => {
      const messageDiv = document.createElement('div');
      messageDiv.id = 'message';
      messageDiv.innerHTML = `Club, which is named ${data.name}, created successfully! Redirecting to club page...`;
      document.body.appendChild(messageDiv);

      setTimeout(() => {
        window.location.href = `/clubs/${data.clubId}/`;
      }, 2000);
    })
    .catch(error => {
      console.error("There has been a problem with your fetch operation:", error);
    });
}
