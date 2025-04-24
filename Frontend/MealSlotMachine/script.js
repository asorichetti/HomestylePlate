const items = {
  protein: ['Chicken', 'Turkey', 'Sausage', 'Steak', 'Roast Beef', 'Pork Chops', 'Schnitzel', 'Pork Tenderloin', 'Salmon', 'Breaded Fish', 'Brisket', 'Ham', 'Flank Steak', 'Meatloaf', 'Mini-Meatloafs', 'Burgers'],
  vegetable: ['Broccoli', 'Carrots', 'Asparagus', 'Peas', 'Corn', 'Mixed Vegetables', 'California Vegetables', 'Broccolini', 'Brussel Sprouts', 'Cauliflower', 'Coleslaw', 'Spinnach Salad', 'Caesar Salad', 'Green Beans', 'Carrots and Beans Medley'],
  starch: ['Rice', 'Brown Rice', 'Gnocchi', 'Mashed Potatoes', 'Baked Potato', 'Roasted Potatoes', 'Sweet Potato Mash', 'Baked Sweet Potato', 'Orzo', 'Stuffing', 'Baked Beans', 'Garlic Bread', 'Fries', 'Sweet Potato Fries', 'Onion Rings', 'Spring Rolls', 'Potstickers']
};

let lockedItems = { protein: false, vegetable: false, starch: false };

function pullLever() {
  const lever = document.getElementById('lever');
  lever.classList.add('lever-pulled');

  setTimeout(() => {
    lever.classList.remove('lever-pulled');
  }, 500);
}

function spinWheel(type) {
  if (lockedItems[type]) return;

  const strip = document.getElementById(`${type}-strip`);
  const itemsArray = items[type];
  const randomIndex = Math.floor(Math.random() * itemsArray.length);
  strip.textContent = itemsArray[randomIndex];
}

function toggleLock(type) {
  lockedItems[type] = !lockedItems[type];
  const button = document.getElementById(`lock-${type}`);
  const icon = lockedItems[type] ? 'fa-lock' : 'fa-lock-open';
  const label = lockedItems[type] ? 'Unlock' : 'Lock';
  button.setAttribute('aria-pressed', lockedItems[type]);
  button.innerHTML = `<i class="fas ${icon}" aria-hidden="true"></i> ${label}`;
}

function handleSpin() {
  pullLever();
  spinWheel('protein');
  spinWheel('vegetable');
  spinWheel('starch');
}