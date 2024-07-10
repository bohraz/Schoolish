const feed = document.getElementById("feed");

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

function loadPosts() {

}

function loadPost() {

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
    }

    feed.append(...posts);

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