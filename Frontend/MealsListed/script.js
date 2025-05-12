document.addEventListener('DOMContentLoaded', () => {
    const gallery = document.getElementById('meal-gallery');
    const form = document.getElementById('meal-form');
    const lockButton = document.querySelector('.lock-btn');
    let requestedHF = 0;
    let requestedP3 = 0;


    //Utility Function to get locked meals from sessionStorage
    function getLockedMeals(){
        const data = sessionStorage.getItem('lockedMeals');
        return data ? JSON.parse(data) : [];
    }
    //Utility Function to save locked meals to sessionStorage
    function saveLockedMeals(meals){
        sessionStorage.setItem('lockedMeals', JSON.stringify(meals));
    }
    

    //Fetch meals from backend on page load
    function fetchMeals(){
        const lockedMeals = getLockedMeals();
        const lockedHF = lockedMeals.filter(m => m.source === 'HF_Meal').length;
        const lockedP3 = lockedMeals.filter(m => m.source === 'P3_Meal').length;

        //Request meals with consideration for locked ones
        const neededHF = Math.max(0, requestedHF - lockedHF);
        const neededP3 = Math.max(0, requestedP3 - lockedP3);

        const apiUrl = `http://localhost:8080/meals?type=HF_Meal:${neededHF},P3_Meal:${neededP3}`;

        fetch(apiUrl)
            .then(response => {
                if (!response.ok){
                    throw new Error(`Failed to fetch meals: ${response.status}`);
                }
                return response.json();
            })
            .then(meals => {
                renderMeals(meals);
            })
            .catch(error => {
                console.error('Error fetching meals:', error);
                gallery.innerHTML = '<p>Error loading meals. Please try again later.</p>';
            });
    }

    //Render meals and include locked ones
    function renderMeals(meals){
        const lockedMeals = getLockedMeals();
        meals.forEach(meal => {
            const mealCard = createMealCard(meal);
            //Check if meal is locked and apply styles accordingly
            if (meal.locked){
                mealCard.querySelector('.lock-btn').classList.add('locked');
                mealCard.querySelector('.lock-btn').textContent = 'Unlock';
            }
            gallery.appendChild(mealCard);
        });

        //Add locked meals on top of the page load
        lockedMeals.forEach(meal => {
            const mealCard = createMealCard(meal);
            mealCard.querySelector('.lock-btn').classList.add('locked');
            mealCard.querySelector('.lock-btn').textContent = 'Unlock';
            gallery.appendChild(mealCard);
        });
    }

    
})