document.addEventListener('DOMContentLoaded', () => {
    const gallery = document.getElementById('meal-gallery');
    const leftArrow = document.getElementById('left-arrow');
    const rightArrow = document.getElementById('right-arrow');
    let currentIndex = 0;
    let meals = [];

    // Fetch meals from backend on page load
    fetch('http://backend:8080/meals?type=HF_Meal:4,P3_Meal:4')
        .then(response => {
            if (!response.ok) {
                throw new Error(`Failed to fetch meals: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            meals = data;
            if (!Array.isArray(meals) || meals.length === 0) {
                gallery.innerHTML = '<p>No meals found.</p>';
                return;
            }

            renderMeal(currentIndex);
        })
        .catch(error => {
            console.error('Error fetching meals:', error);
            gallery.innerHTML = '<p>Error loading meals. Please try again later.</p>';
        });

    // Render meal image and link based on current index
    function renderMeal(index) {
        const meal = meals[index];
        gallery.innerHTML = '';

        const card = document.createElement('div');
        card.className = 'meal-card';

        const img = document.createElement('img');
        img.src = meal.image_url;
        img.alt = meal.name;
        img.className = 'meal-image';

        const nameLink = document.createElement('a');
        nameLink.href = meal.recipe_link;
        nameLink.textContent = meal.name;
        nameLink.className = 'meal-link';
        nameLink.target = '_blank';
        nameLink.rel = 'noopener';

        card.appendChild(img);
        card.appendChild(nameLink);
        gallery.appendChild(card);
    }

    // Move to the next meal
    function nextMeal() {
        currentIndex = (currentIndex + 1) % meals.length;
        renderMeal(currentIndex);
    }

    // Move to the previous meal
    function prevMeal() {
        currentIndex = (currentIndex - 1 + meals.length) % meals.length;
        renderMeal(currentIndex);
    }

    // Event listeners for arrow buttons
    leftArrow.addEventListener('click', prevMeal);
    rightArrow.addEventListener('click', nextMeal);
});
