document.addEventListener('DOMContentLoaded', () => {
    const gallery = document.getElementById('meals-list');
    const reloadButton = document.getElementById('reload-meals');
    const urlParams = new URLSearchParams(window.location.search);

    const shouldReset = urlParams.get('reset') === 'true';
    if (shouldReset) {
        console.log('Resetting locked meals from sessionStorage');
        sessionStorage.removeItem('lockedMeals');
    }

    const requestedHF = parseInt(urlParams.get('hf')) || 0;
    const requestedP3 = parseInt(urlParams.get('p3')) || 0;

    function getLockedMeals() {
        const raw = sessionStorage.getItem('lockedMeals');
        if (!raw) return [];
        try {
            return JSON.parse(raw);
        } catch {
            return [];
        }
    }

    function saveLockedMeals(locked) {
        sessionStorage.setItem('lockedMeals', JSON.stringify(locked));
    }

    function createMealCard(meal, isLocked) {
        const card = document.createElement('div');
        card.classList.add('meal-entry');

        const img = document.createElement('img');
        img.src = meal.image_url;
        img.alt = meal.name;
        img.classList.add('meal-image');

        const content = document.createElement('div');
        content.classList.add('meal-content');

        const title = document.createElement('div');
        title.classList.add('meal-title');
        title.textContent = meal.name;

        const buttons = document.createElement('div');
        buttons.classList.add('meal-buttons');

        const recipeLink = document.createElement('a');
        recipeLink.href = meal.recipe_link;
        recipeLink.textContent = 'View Recipe';
        recipeLink.target = '_blank';
        recipeLink.rel = 'noopener noreferrer';
        recipeLink.classList.add('view-recipe-btn');

        const lockBtn = document.createElement('button');
        lockBtn.classList.add('lock-btn');
        lockBtn.textContent = isLocked ? 'Unlock' : 'Lock';
        if (isLocked) lockBtn.classList.add('locked');

        buttons.appendChild(recipeLink);
        buttons.appendChild(lockBtn);
        content.appendChild(title);
        content.appendChild(buttons);
        card.appendChild(img);
        card.appendChild(content);

        lockBtn.addEventListener('click', () => {
            const lockedMeals = getLockedMeals();
            if (isLocked) {
                const updated = lockedMeals.filter(m => m.name !== meal.name);
                saveLockedMeals(updated);
                
            } else {
                if (!lockedMeals.some(m => m.name === meal.name)) {
                    lockedMeals.push({
                        name: meal.name,
                        rating: meal.rating,
                        image_url: meal.image_url,
                        recipe_link: meal.recipe_link,
                        source: meal.source
                    });
                    saveLockedMeals(lockedMeals);
                    
                }
            }
        });

        return card;
    }

    function renderMeals(fetchedMeals, lockedMeals) {
        gallery.innerHTML = '';

        lockedMeals.forEach(meal => {
            const card = createMealCard(meal, true);
            gallery.appendChild(card);
        });

        fetchedMeals.forEach(meal => {
            const card = createMealCard(meal, false);
            gallery.appendChild(card);
        });
    }

    function fetchMeals() {
        const lockedMeals = getLockedMeals();
        const lockedHF = lockedMeals.filter(m => m.source === 'HF_Meal').length;
        const lockedP3 = lockedMeals.filter(m => m.source === 'P3_Meal').length;

        const neededHF = Math.max(0, requestedHF - lockedHF);
        const neededP3 = Math.max(0, requestedP3 - lockedP3);

        const queryParts = [];
        if (neededHF > 0) queryParts.push(`HF_Meal:${neededHF}`);
        if (neededP3 > 0) queryParts.push(`P3_Meal:${neededP3}`);

        if (queryParts.length === 0) {
            renderMeals([], lockedMeals);
            return;
        }

        const apiUrl = `/meals?type=${queryParts.join(',')}`;

        fetch(apiUrl)
            .then(res => {
                if (!res.ok) throw new Error(`Failed to fetch meals: ${res.status}`);
                return res.json();
            })
            .then(fetchedMeals => {
                renderMeals(fetchedMeals, lockedMeals);
            })
            .catch(err => {
                console.error('Error fetching meals:', err);
                gallery.innerHTML = '<p>Error loading meals. Please try again later.</p>';
            });
    }

    reloadButton?.addEventListener('click', fetchMeals);

    fetchMeals();
});
