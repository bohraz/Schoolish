const postId = window.location.pathname.split("/")[2];
const wsUrl = "ws://127.0.0.1:80/ws/";
const sockets = [wsUrl + "post?id=" + postId, wsUrl + "comment?id=" + postId];

sockets.forEach(v => {
    const socket = new WebSocket(v);

    // Connection opened
    socket.addEventListener("open", function (event) {
        console.log("Connected to WebSocket server");
    });

    // Listen for messages
    socket.addEventListener("message", function (event) {
        console.log("Message from server ", event.data);
    });

    // Connection closed
    socket.addEventListener("close", function (event) {
        console.log("Disconnected from WebSocket server");
    });
});


function submitComment(event) {
    event.preventDefault();

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
        totalComments = post.comments;
    } catch (error) {
        console.error("Error:", error);
    }

    fetchComments();
}

async function fetchComments() {
    try {
        const getCommentsUrl = `/api/comment?postId=${window.location.pathname.split("/")[2]}&limit=${limit}&offset=${offset}`;
        const response = await fetch(getCommentsUrl);
        if (!response.ok) {
            throw new Error("Network response was not ok");
        }

        offset += limit;

        const comments = await response.json();

        comments.forEach(comment => {
            const commentDiv = document.createElement('div');
            commentDiv.className = 'comment';
            commentDiv.id = comment.id;
            commentDiv.innerHTML = `
                <h4 class="commentAuthor">${comment.user.firstName + " " + comment.user.lastName}</h4>
                <p class="commentContent">${comment.content}</p>
            `;
            document.getElementById('comments').appendChild(commentDiv);
        });
    } catch (error) {
        console.error("Error:", error);
    }

    requestAnimationFrame(() => {
        if (offset < totalComments && window.innerHeight >= document.body.scrollHeight) {
            fetchComments();
        }
    })
}

let title = document.getElementById("title");
let author = document.getElementById("author");
let content = document.getElementById("content");

const limit = 5;
let offset = 0;
let totalComments = 0;

const postUrl = `/api/post?id=${window.location.pathname.split("/")[2]}`;
const commentUrl = "/api/comment/create/";

fetchPost();

function throttle(fn, wait) {
    let lastTime = 0;
    return function (...args) {
        const now = new Date().getTime();
        if (now - lastTime >= wait) {
            lastTime = now;
            fn.apply(this, args);
        }
    };
}

window.addEventListener("scroll", throttle(function () {
    if (offset < totalComments && window.scrollY + window.innerHeight >= document.body.scrollHeight) {
        fetchComments();
    }
}, 25));