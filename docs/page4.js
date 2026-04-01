// Image Map Manager - Handles responsive coordinates
class ImageMapManager {
    constructor(imageId, mapName, originalWidth, originalHeight) {
        this.image = document.getElementById(imageId);
        this.map = document.querySelector(`map[name="${mapName}"]`);
        this.originalWidth = originalWidth;
        this.originalHeight = originalHeight;
        
        if (!this.image) {
            console.error(`Image "${imageId}" not found`);
            return;
        }
        
        if (!this.map) {
            console.error(`Map "${mapName}" not found`);
            return;
        }
        
        console.log('✅ Image Map initialized');
        console.log('Image:', this.image.src);
        console.log('Map areas:', this.map.querySelectorAll('area').length);
        
        // Store original coordinates
        this.areas = Array.from(this.map.querySelectorAll('area'));
        this.areas.forEach(area => {
            area.originalCoords = area.getAttribute('coords');
            console.log('Original coordinates:', area.originalCoords);
        });
        
        // Handle resize
        this.handleResize = this.handleResize.bind(this);
        window.addEventListener('resize', this.handleResize);
        
        // Initial resize after image loads
        if (this.image.complete) {
            this.handleResize();
        } else {
            this.image.addEventListener('load', () => this.handleResize());
        }
        
        // Add click debugging
        this.addClickTracking();
    }
    
    handleResize() {
        if (!this.image) return;
        
        const currentWidth = this.image.clientWidth;
        const currentHeight = this.image.clientHeight;
        
        if (currentWidth === 0 || currentHeight === 0) return;
        
        const scaleX = currentWidth / this.originalWidth;
        const scaleY = currentHeight / this.originalHeight;
        
        console.log(`Resizing map: ${currentWidth}x${currentHeight}, scale: ${scaleX.toFixed(3)}x${scaleY.toFixed(3)}`);
        
        this.areas.forEach(area => {
            if (area.originalCoords) {
                const coords = area.originalCoords.split(',').map((coord, index) => {
                    const value = parseInt(coord);
                    return Math.round(value * (index % 2 === 0 ? scaleX : scaleY));
                }).join(',');
                area.setAttribute('coords', coords);
            }
        });
    }
    
    addClickTracking() {
        this.areas.forEach(area => {
            area.addEventListener('click', (e) => {
                console.log('🎯 Image Map Area Clicked!');
                console.log('  Target:', area.getAttribute('alt'));
                console.log('  Navigating to:', area.getAttribute('href'));
            });
        });
    }
}

// Auto-hide controls
class AutoHideControls {
    constructor(selector, delay = 3000) {
        this.controls = document.querySelector(selector);
        this.delay = delay;
        
        if (this.controls) {
            this.init();
        }
    }
    
    show() {
        this.controls.style.opacity = '0.7';
        clearTimeout(this.timeout);
        this.timeout = setTimeout(() => this.hide(), this.delay);
    }
    
    hide() {
        this.controls.style.opacity = '0.2';
    }
    
    init() {
        this.show();
        this.controls.addEventListener('mousemove', () => this.show());
        this.controls.addEventListener('click', () => this.show());
    }
}

// Initialize everything
document.addEventListener('DOMContentLoaded', () => {
    console.log('Page4 loading with Image Map approach...');
    
    // Initialize image map (original SVG: 500x400)
    const imageMap = new ImageMapManager('bankingImage', 'bankingMap', 500, 400);
    
    // Initialize auto-hide controls
    const controls = new AutoHideControls('.back-button-container', 3000);
    
    // Add visual hover effect for image map areas
    const img = document.getElementById('bankingImage');
    const areas = document.querySelectorAll('area');
    
    areas.forEach(area => {
        area.addEventListener('mouseenter', () => {
            if (img) img.style.filter = 'brightness(0.97)';
        });
        area.addEventListener('mouseleave', () => {
            if (img) img.style.filter = 'brightness(1)';
        });
    });
    
    // INVISIBLE BUTTON TRACKING
    const statementButton = document.getElementById('statementButton');
    if (statementButton) {
        console.log('✅ Statement button ready - Click the bottom-left corner to navigate to page5');
        statementButton.addEventListener('click', (e) => {
            console.log('📄 Statement button clicked! Navigating to page5...');
            // The href will handle the navigation automatically
        });
    }
    
    console.log('✅ Page4 ready - Click on the SELECT CREDIT CARD button or the invisible STATEMENT button in bottom-left corner');
});