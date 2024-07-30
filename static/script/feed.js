const feed = document.getElementById("feed");
var unloadedPosts = [];
var looping = false;

const socket = new WebSocket("ws://127.0.0.1:80/ws/feed")

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

function createPost(event) {
    event.preventDefault();

    const url = "/api/post/create/";
    const data = {
        title: document.getElementById("title").value,
        content: document.getElementById("content").value
    }

    fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
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
                console.log("Post created successfully!");
                setTimeout(() => {
                    window.location.href = `/post/${data.id}`;
                }, 3000);
            } else {
                console.error("Failed to create post!");
            }
        })
        .catch(error => {
            console.error("Error:", error);
        })
}

async function getPosts() {
    if (looping) {
        return;
    }

    looping = true;

    while (unloadedPosts.length > 0) {
        try {
            const url = "/api/posts?amount=3"

            const response = await fetch(url);
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }

            const posts = await response.json()

            for (let post of posts) {
                const postDiv = unloadedPosts.shift();
                postDiv.querySelector("h3").textContent = post.title;
                postDiv.querySelector("h4").textContent = post.author;
                postDiv.querySelector("p").textContent = post.content;
            }
        } catch (error) {
            console.error("Error:", error);
        }
    }

    looping = false;
}

function generateBlankPosts() {
    let posts = [];
    for (let i = 0; i < 3; i++) {
        const post = document.createElement("div");
        post.classList.add("post");
        post.innerHTML = `
            <h3>Post Title</h3>
            <h4>Post Author</h4>
            <p>Post Content</p>
        `

        posts.push(post);
        unloadedPosts.push(post);
    }

    feed.append(...posts);

    getPosts();

    requestAnimationFrame(() => {
        if (window.innerHeight >= document.body.scrollHeight) {
            generateBlankPosts();
        }
    })
}

generateBlankPosts();

window.addEventListener("scroll", function () {
    if (window.scrollY + window.innerHeight >= document.body.scrollHeight) {
        generateBlankPosts();
    }
});