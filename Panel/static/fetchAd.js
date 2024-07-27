(function () {
  const publisherID = document.currentScript.getAttribute('id');
  console.log("Meoooooooooooooo" + publisherID);
  const adContainer = document.getElementById('adBox');
  console.log("SAGGGGGGGGGGGGG" + adContainer);
  if (!adContainer) {
    console.error('Ad container element not found.');
    return;
  }
  function fetchAd() {
    fetch(`https://adserver.lontra.tech/api/ads?publisherID=${publisherID}`)
      .then(response => response.json())
      .then(data => {
        console.log("Received data:", JSON.stringify(data, null, 2));
        if (data && data.length > 0) {
          const ad = data;
          const adContent = `
                <img src="${ad.ImagePath}" alt="${ad.Title}" style="width:100%;" />
                <h3>${ad.Title}</h3>
                <a ${ad.ClickLink} target="_blank" class="click-here">Click here</a>
              `;
              // delete href
          adContainer.innerHTML = adContent;
          const observer = new IntersectionObserver((entries) => {
            if (entries[0].isIntersecting) {
              fetch(ad.impressionLink);
              observer.disconnect();
            }
          }, { threshold: 1.0 });
          observer.observe(adContainer);
          document.querySelector('.click-here').addEventListener('click', function (event) {
            event.preventDefault();
            fetch(ad.ClickLink, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json'
              },
              body: JSON.stringify({ ad: ad.Title, clickedAt: new Date().toISOString() })
            })
              .then(response => response.json())
              .then(data => {
                console.log('Click recorded:', data);
                // window.open(ad.ClickLink, '_blank');
                // dont use it
              })
              .catch(error => {
                console.error('Error recording click:', error);
              });
          });
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