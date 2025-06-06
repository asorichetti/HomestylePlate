document.addEventListener('DOMContentLoaded', () => {
    const gallery = document.getElementById('meal-gallery');
    const leftArrow = document.getElementById('left-arrow');
    const rightArrow = document.getElementById('right-arrow');
    let currentIndex = 0;
    let meals = [];

    // ----- FORM SUBMISSION HANDLER -----
    const form = document.getElementById('meal-form');
    if (form) {
        form.addEventListener('submit', (event) => {
            event.preventDefault(); // Prevent normal form submission

            // Parse inputs for HF and P3 meal counts
            const hfInput = form.querySelector('input[name="hf"]');
            const p3Input = form.querySelector('input[name="p3"]');
            const hfCount = hfInput ? parseInt(hfInput.value) || 0 : 0;
            const p3Count = p3Input ? parseInt(p3Input.value) || 0 : 0;

            // Clear sessionStorage to reset locked meals for new session
            sessionStorage.removeItem('lockedMeals');

            // Redirect to MealsListed with counts and reset=true
            const query = new URLSearchParams({
                hf: hfCount,
                p3: p3Count,
                reset: 'true'
            });

            const redirectURL = `${form.action}?${query.toString()}`;
            window.location.href = redirectURL;
        });
    }

    // ----- FETCH INITIAL GALLERY MEALS -----
    fetch('/meals?type=HF_Meal:4,P3_Meal:4')
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

    // ----- RENDER A SINGLE MEAL -----
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

    // ----- NAVIGATION ARROWS -----
    function nextMeal() {
        currentIndex = (currentIndex + 1) % meals.length;
        renderMeal(currentIndex);
    }

    function prevMeal() {
        currentIndex = (currentIndex - 1 + meals.length) % meals.length;
        renderMeal(currentIndex);
    }

    leftArrow?.addEventListener('click', prevMeal);
    rightArrow?.addEventListener('click', nextMeal);
});

