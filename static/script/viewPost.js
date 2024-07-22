let title = document.getElementById("title");
let author = document.getElementById("author");
let content = document.getElementById("content");

const postUrl = `/api/post?id=${window.location.pathname.split("/")[2]}`;

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

fetchPost();