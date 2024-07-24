function login(event) {
    event.preventDefault();

    const url = "/api/login/";
    const data = {
        username: document.getElementById("username").value,
        password: document.getElementById("password").value
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
            if (data.success) {
                console.log('Logged in successfully!');
            } else {
                console.error('Failed to log in!');
            }
        })
        .catch(error => {
            console.error('Error:', error);
        })
}