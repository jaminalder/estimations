(function() {
  // Theme switcher
  const root = document.documentElement;
  const dropdown = document.getElementById('themeDropdown');
  const button = document.getElementById('themeButton');
  const label = document.getElementById('themeLabel');

  function applyTheme(theme) {
    root.classList.remove('theme-light', 'theme-dark');
    if (theme === 'light') root.classList.add('theme-light');
    if (theme === 'dark') root.classList.add('theme-dark');
  }

  function setTheme(theme) {
    try { localStorage.setItem('theme', theme); } catch (_) {}
    applyTheme(theme);
    if (label) { label.textContent = theme.charAt(0).toUpperCase() + theme.slice(1); }
  }

  const saved = (function() { try { return localStorage.getItem('theme') || 'system'; } catch (_) { return 'system'; } })();
  setTheme(saved);

  if (button && dropdown) {
    button.addEventListener('click', function(e) {
      e.stopPropagation();
      dropdown.classList.toggle('is-active');
    });
    document.addEventListener('click', function() {
      dropdown.classList.remove('is-active');
    });
    dropdown.querySelectorAll('[data-theme]').forEach(function(item) {
      item.addEventListener('click', function(e) {
        e.preventDefault();
        e.stopPropagation();
        const theme = this.getAttribute('data-theme');
        setTheme(theme);
        dropdown.classList.remove('is-active');
      });
    });
  }

  // Room card selection
  const deck = document.querySelector('.card-deck');
  const cards = deck ? Array.from(deck.querySelectorAll('.poker-card')) : [];
  const voteBtn = document.getElementById('voteButton');
  let selectedCard = null;

  function resetCards() {
    cards.forEach(function(c) {
      c.classList.remove('has-background-primary', 'has-text-white');
      c.classList.add('has-background-white', 'has-text-dark');
    });
  }

  if (cards.length) {
    cards.forEach(function(card) {
      card.addEventListener('click', function() {
        resetCards();
        this.classList.remove('has-background-white', 'has-text-dark');
        this.classList.add('has-background-primary', 'has-text-white');
        selectedCard = (this.textContent || '').trim();
        if (voteBtn) voteBtn.disabled = false;
      });
    });
  }

  if (voteBtn) {
    voteBtn.addEventListener('click', function() {
      if (!selectedCard) return;
      this.disabled = true;
      this.classList.remove('is-success');
      this.classList.add('is-primary');
      this.innerHTML = '<span class="icon"><i class="fas fa-check"></i></span><span>Voted!</span>';
    });
  }
})();

