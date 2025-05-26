const GALLERY = document.getElementById('dessertGallery');
const MODAL = document.getElementById('dessertModal');
const MODAL_TITLE = document.getElementById('modalTitle');
const MODAL_IMAGE = document.getElementById('modalImage');
const MODAL_INGREDIENTS = document.getElementById('modalIngredients');
const MODAL_BAKE_TIME = document.getElementById('modalBakeTime');
const MODAL_LINK = document.getElementById('modalRecipeLink');
const CLOSE_BTN = document.querySelector(".close-btn");
const PREVPAGEBTN = document.getElementById("prev-page");
const NEXTPAGEBTN = document.getElementById("next-page");
const PAGEINFO = document.getElementById("page-number")

let currentPage = 1;
const dessertsPerPage = 10;
let desserts = [];

async function fetchDesserts(){
    try{
        const res = await fetch("http://localhost:8081/api/desserts");
        const data = await res.json();
        desserts = data;
        renderPage(currentPage);
    } catch (err) {
        console.error("Failed to Fetch Desserts:", err);
    }
}

function renderPage(page){
    GALLERY.innerHTML = "";
    const start = (page - 1) * dessertsPerPage;
    const end = start + dessertsPerPage;
    const pageDesserts = desserts.slice(start, end);

    pageDesserts.forEach(dessert => {
        const card = document.createElement("div");
        card.className = "polaroid";
        card.innerHTML = `
            <img src="${dessert.image_url}" alt="${dessert.name}" />
            <h3>${dessert.name}</h3>
            `;
        card.addEventListener("click", () => openModal(dessert));
        GALLERY.appendChild(card);
    });

    PAGEINFO.textContent = `Page ${currentPage} of ${Math.ceil(desserts.length / dessertsPerPage)}`;
    PREVPAGEBTN.disabled = currentPage === 1;
    NEXTPAGEBTN.disabled = currentPage === Math.ceil(desserts.length/dessertsPerPage);
}

function openModal(dessert){
    MODAL_TITLE.textContent = dessert.name;
    MODAL_IMAGE.src = dessert.image_url;
    MODAL_IMAGE.alt = dessert.name
    MODAL_INGREDIENTS.textContent = dessert.ingredients;
    MODAL_BAKE_TIME.textContent = dessert.bake_time;
    MODAL_LINK.href = dessert.recipe_link;
    MODAL.classList.remove("hidden");
}

CLOSE_BTN.addEventListener("click", () => {
    MODAL.classList.add("hidden");
});

PREVPAGEBTN.addEventListener("click", () => {
    if (currentPage > 1){
        currentPage--;
        renderPage(currentPage);
    }
});

NEXTPAGEBTN.addEventListener("click", () => {
    if (currentPage < Math.ceil(desserts.length/dessertsPerPage)){
        currentPage++;
        renderPage(currentPage);
    }
});

fetchDesserts();