(function () {
  const publisherID = document.currentScript.getAttribute('id');
  const adContainer = document.getElementById('adBox');
  if (!adContainer) {
    console.error('Ad container element not found.');
    return;
  }
  function fetchAd() {
    fetch(`https://adserver.lontra.tech/api/ads?publisherID=${publisherID}`)
      .then(response => response.json())
      .then(data => {
        console.log("Received data:", JSON.stringify(data, null, 2));
        if (data) {
          const ad = data;
          const adContent = `
                <img src="https://panel.lontra.tech/${ad.ImagePath}" alt="${ad.Title}" style="width:100%;" />
                <h3>${ad.Title}</h3>
                <a target="_blank" class="click-here" href="${ad.ClickLink}">Click here</a>
              `;
              // delete href
          adContainer.innerHTML = adContent;
          const observer = new IntersectionObserver((entries) => {
            if (entries[0].isIntersecting) {
              fetch(ad.ImpressionLink);
              observer.disconnect();
            }
          }, { threshold: 1.0 });
          observer.observe(adContainer);
        } else {
          adContainer.innerHTML = '<div class="ad-content">No ad available</div>';
        }
      })
      .catch(error => {
        console.error('Error fetching ad:', error);
        adContainer.innerHTML = '<div class="ad-content">Error loading ad</div>';
      });
  }
  function placeAd() {
    const paragraphs = document.getElementsByTagName('p');
    if (paragraphs.length > 0) {
      const lastParagraph = paragraphs[paragraphs.length - 1];
      lastParagraph.parentNode.insertBefore(adContainer, lastParagraph.nextSibling);
    } else {
      document.body.appendChild(adContainer);
    }
  }
  fetchAd();
  placeAd();
})();