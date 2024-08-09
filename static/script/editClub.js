function edit(event) {
    event.preventDefault();

    const path = window.location.pathname;
    const clubId = path.split("/")[2];
    const url = `/api/club/${clubId}/edit/`;

    const data = {
        id: clubId,
        name: document.getElementById("name").value,
        description: document.getElementById("description").value,
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
                console.log('Club edited successfully!');
            } else {
                console.error('Failed to edit club!');
            }
        })
        .catch(error => {
            console.error('Error:', error);
        })
}