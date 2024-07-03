function rando() {
  console.log("This works!");
  // Sending a POST request to the Go server
  // fetch("http://localhost:80/api/", {
  //   method: "POST", // Specify the method
  //   headers: {
  //     "Content-Type": "application/json", // Set the content type header
  //   },
  //   body: JSON.stringify({
  //     key: "value", // Your request body content
  //   }),
  // })
  //   .then((response) => {
  //     if (!response.ok) {
  //       throw new Error("Network response was not ok");
  //     }
  //     return response.json(); // or .text() if the response is not JSON
  //   })
  //   .then((data) => {
  //     console.log(data); // Process the data
  //   })
  //   .catch((error) => {
  //     console.error("There was a problem with your fetch operation:", error);
  //   });
}

console.log("Loaded javascript!")
