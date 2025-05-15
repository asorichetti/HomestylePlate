document.addEventListener('DOMContentLoaded', () => {
    const menuToggle = document.querySelector('.menu-toggle');
    const navLinks = document.querySelector('.nav-links');
  
    if (!menuToggle || !navLinks) return;
  
    // Function to toggle menu
    function toggleMenu() {
      const isActive = navLinks.classList.toggle('active');
      menuToggle.classList.toggle('active');
      menuToggle.setAttribute('aria-expanded', isActive);
    }
  
    // Close menu on link click (only if menu is open)
    function closeMenuOnLinkClick(e) {
      if (window.innerWidth <= 768 && navLinks.classList.contains('active')) {
        toggleMenu();
      }
    }
  
    // Event listeners
    menuToggle.addEventListener('click', toggleMenu);
  
    navLinks.querySelectorAll('a').forEach(link => {
      link.addEventListener('click', closeMenuOnLinkClick);
    });
  
    // Optional: Close menu if window resized larger than breakpoint
    window.addEventListener('resize', () => {
      if (window.innerWidth > 768 && navLinks.classList.contains('active')) {
        toggleMenu();
      }
    });
  });
  