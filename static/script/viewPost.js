function submitComment(event) {
    console.log("WHY WONT WORK!???")
    event.preventDefault();
    console.log("Submitting comment!");

    const sendData = {
        postId: Number(window.location.pathname.split("/")[2]),
        content: document.getElementById("contentInput").value
    };

    console.log(typeof (sendData.postId), sendData.postId)

    fetch(commentUrl, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(sendData)
    })
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }
            return response.json();
        })
        .then(data => {
            const commentDiv = document.createElement('div');
            commentDiv.className = 'comment';
            commentDiv.id = data.id;
            commentDiv.innerHTML = `
                <h4 class="commentAuthor">${data.user.firstName + " " + data.user.lastName}</h4>
                <p class="commentContent">${data.content}</p>
            `;
            document.getElementById('comments').appendChild(commentDiv);
        })
        .catch(error => {
            console.error("Error:", error);
        });
}

async function fetchPost() {
    try {
        const response = await fetch(postUrl);
        if (!response.ok) {
            throw new Error("Network response was not ok");
        }

        const post = await response.json();

        title.textContent = post.title;
        author.textContent = post.userId;
        content.textContent = post.content;
    } catch (error) {
        console.error("Error:", error);
    }
}

let title = document.getElementById("title");
let author = document.getElementById("author");
let content = document.getElementById("content");

const postUrl = `/api/post?id=${window.location.pathname.split("/")[2]}`;
const commentUrl = "/api/comment/create/";

fetchPost();