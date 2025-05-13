document.addEventListener('DOMContentLoaded', () => {
    const gallery = document.getElementById('meals-list');
    const reloadButton = document.getElementById('reload-meals');
    const urlParams = new URLSearchParams(window.location.search);
    let requestedHF = parseInt(urlParams.get('hf')) || 0;
    let requestedP3 = parseInt(urlParams.get('p3')) || 0;

    // Utility function to get locked meals from sessionStorage
    function getLockedMeals() {
        const data = sessionStorage.getItem('lockedMeals');
        return data ? JSON.parse(data) : [];
    }

    // Utility function to save locked meals to sessionStorage
    function saveLockedMeals(meals) {
        sessionStorage.setItem('lockedMeals', JSON.stringify(meals));
    }

    // Fetch meals from backend on page load
    function fetchMeals() {
        const lockedMeals = getLockedMeals();
        const lockedHF = lockedMeals.filter(m => m.source === 'HF_Meal').length;
        const lockedP3 = lockedMeals.filter(m => m.source === 'P3_Meal').length;

        // Adjust the number of meals requested based on the locked ones
        const neededHF = Math.max(0, requestedHF - lockedHF);
        const neededP3 = Math.max(0, requestedP3 - lockedP3);

        // API URL for meal request
        const apiUrl = `http://localhost:8080/meals?type=HF_Meal:${neededHF},P3_Meal:${neededP3}`;

        fetch(apiUrl)
            .then(response => {
                if (!response.ok) {
                    throw new Error(`Failed to fetch meals: ${response.status}`);
                }
                return response.json();
            })
            .then(meals => {
                renderMeals(meals, lockedMeals);
            })
            .catch(error => {
                console.error('Error fetching meals:', error);
                gallery.innerHTML = '<p>Error loading meals. Please try again later.</p>';
            });
    }

    // Render meals and include locked ones
    function renderMeals(meals, lockedMeals) {
        gallery.innerHTML = ''; // Clear any previous meals displayed

        // Render locked meals
        lockedMeals.forEach(meal => {
            const mealCard = createMealCard(meal);
            mealCard.querySelector('.lock-btn').classList.add('locked');
            mealCard.querySelector('.lock-btn').textContent = 'Unlock';
            gallery.appendChild(mealCard);
        });

        // Render available meals
        meals.forEach(meal => {
            const mealCard = createMealCard(meal);
            gallery.appendChild(mealCard);
        });
    }

    // Create a meal card element
    function createMealCard(meal) {
        const card = document.createElement('div');
        card.classList.add('meal-entry');
    
        // Image on the left
        const img = document.createElement('img');
        img.src = meal.image_url;
        img.alt = meal.name;
        img.classList.add('meal-image');
    
        // Right content block
        const content = document.createElement('div');
        content.classList.add('meal-content');
    
        // Meal title
        const title = document.createElement('div');
        title.classList.add('meal-title');
        title.textContent = meal.name;
    
        // Buttons container
        const buttonGroup = document.createElement('div');
        buttonGroup.classList.add('meal-buttons');
    
        // View Recipe button (styled link)
        const link = document.createElement('a');
        link.href = meal.recipe_link;
        link.textContent = 'View Recipe';
        link.target = '_blank';
        link.rel = 'noopener';
        link.classList.add('view-recipe-btn');
    
        // Lock/Unlock button
        const lockBtn = document.createElement('button');
        lockBtn.classList.add('lock-btn');
        lockBtn.textContent = 'Lock';
    
        // Append buttons
        buttonGroup.appendChild(link);
        buttonGroup.appendChild(lockBtn);
    
        // Assemble content
        content.appendChild(title);
        content.appendChild(buttonGroup);
    
        // Assemble card
        card.appendChild(img);
        card.appendChild(content);
    
        // Lock toggle logic
        lockBtn.addEventListener('click', () => {
            meal.locked = !meal.locked;
            lockBtn.classList.toggle('locked', meal.locked);
            lockBtn.textContent = meal.locked ? 'Unlock' : 'Lock';
    
            const currentLocked = getLockedMeals();
            if (meal.locked) {
                meal.source = meal.source || (meal.name.includes('HF') ? 'HF_Meal' : 'P3_Meal');
                currentLocked.push(meal);
            } else {
                const index = currentLocked.findIndex(m => m.name === meal.name);
                if (index !== -1) currentLocked.splice(index, 1);
            }
            saveLockedMeals(currentLocked);
        });
    
        return card;
    }
    

    // Event listener for the reload button
    reloadButton.addEventListener('click', () => {
        gallery.innerHTML = '';
        fetchMeals();
    });

    // Initial fetch when the page loads
    fetchMeals();
});
