const items = {
    protein: ['Chicken', 'Turkey', 'Sausage', 'Steak', 'Roast Beef', 'Pork Chops', 'Schnitzel', 'Salmon'],
    vegetable: ['Broccoli', 'Carrots', 'Asparagus', 'Peas', 'Corn', 'Cauliflower', 'Salad', 'Green Beans'],
    starch: ['Rice', 'Gnocchi', 'Mashed Potatoes', 'Fries', 'Sweet Potato Mash', 'Stuffing', 'Orzo', 'Garlic Bread']
  };
  
  let lockedItems = {
    protein: false,
    vegetable: false,
    starch: false
  };
  
  function toggleLock(type) {
    lockedItems[type] = !lockedItems[type];
    const lockButton = document.getElementById(`lock-${type}`);
    lockButton.setAttribute('aria-pressed', lockedItems[type]);
    lockButton.innerHTML = lockedItems[type]
      ? `<i class="fas fa-lock" aria-hidden="true"></i> Locked`
      : `<i class="fas fa-lock-open" aria-hidden="true"></i> Lock`;
  }
  
  function handleSpin() {
    pullLever();
  }
  
  function pullLever() {
    ['protein', 'vegetable', 'starch'].forEach(type => {
      if (!lockedItems[type]) spinWheel(type);
    });
  
    // Animate lever
    const lever = document.getElementById('lever');
    lever.style.transform = 'translateY(5px)';
    setTimeout(() => (lever.style.transform = 'translateY(0)'), 300);
  }
  
  function spinWheel(type) {
    const reel = document.getElementById(`${type}-reel`);
    const strip = document.getElementById(`${type}-strip`);
  
    strip.innerHTML = '';
    const spinItems = [...items[type], ...items[type]]; // Loop
    spinItems.forEach(text => {
      const item = document.createElement('div');
      item.textContent = text;
      item.style.padding = '5px';
      strip.appendChild(item);
    });
  
    strip.style.transition = 'none';
    strip.style.transform = 'translateY(0)';
    void strip.offsetWidth; // Force reflow
  
    const totalHeight = strip.scrollHeight;
    const spinDistance = -Math.floor(Math.random() * (strip.children.length - 3)) * 30;
  
    strip.style.transition = 'transform 2s cubic-bezier(0.33, 1, 0.68, 1)';
    strip.style.transform = `translateY(${spinDistance}px)`;
  }
  