// Responsive Image Map Manager
class ResponsiveImageMap {
    constructor(imageId, mapName, originalWidth, originalHeight) {
        this.image = document.getElementById(imageId);
        this.map = document.querySelector(`map[name="${mapName}"]`);
        this.originalWidth = originalWidth;
        this.originalHeight = originalHeight;
        this.areas = [];
        
        // Store original coordinates
        if (this.map) {
            this.areas = Array.from(this.map.querySelectorAll('area'));
            this.areas.forEach(area => {
                area.originalCoords = area.getAttribute('coords');
            });
        }
        
        // Bind resize event
        this.handleResize = this.handleResize.bind(this);
        window.addEventListener('resize', this.handleResize);
        
        // Initial calculation
        this.handleResize();
        
        // Add click tracking
        this.addClickTracking();
    }
    
    handleResize() {
        if (!this.image || !this.map) return;
        
        const currentWidth = this.image.clientWidth;
        const currentHeight = this.image.clientHeight;
        
        const scaleX = currentWidth / this.originalWidth;
        const scaleY = currentHeight / this.originalHeight;
        
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
                const href = area.getAttribute('href');
                const alt = area.getAttribute('alt');
                console.log(`Navigation: ${alt} -> ${href}`);
            });
        });
    }
    
    destroy() {
        window.removeEventListener('resize', this.handleResize);
    }
}

// Auto-hide controls manager
class AutoHideControls {
    constructor(selector, delay = 3000) {
        this.controls = document.querySelector(selector);
        this.delay = delay;
        this.timeout = null;
        
        if (this.controls) {
            this.init();
        }
    }
    
    show() {
        if (this.controls) {
            this.controls.style.opacity = '0.7';
            clearTimeout(this.timeout);
            this.timeout = setTimeout(() => this.hide(), this.delay);
        }
    }
    
    hide() {
        if (this.controls) {
            this.controls.style.opacity = '0.2';
        }
    }
    
    init() {
        this.show();
        this.controls.addEventListener('mousemove', () => this.show());
        this.controls.addEventListener('click', () => this.show());
    }
}

// Initialize everything when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    // Initialize responsive image map
    // Original SVG dimensions: 500 x 400
    const imageMap = new ResponsiveImageMap('bankingImage', 'bankingMap', 500, 400);
    
    // Initialize auto-hide controls
    const controls = new AutoHideControls('.back-button-container', 3000);
    
    // Add visual feedback on area hover
    const areas = document.querySelectorAll('area');
    const svgImage = document.getElementById('bankingImage');
    
    areas.forEach(area => {
        area.addEventListener('mouseenter', () => {
            svgImage.style.filter = 'brightness(0.95)';
            svgImage.style.transition = 'filter 0.2s ease';
        });
        
        area.addEventListener('mouseleave', () => {
            svgImage.style.filter = 'brightness(1)';
        });
    });
    
    // Log initialization
    console.log('Page4 initialized with responsive image map');
});

// Handle page unload cleanup
window.addEventListener('beforeunload', () => {
    if (window.imageMap) {
        window.imageMap.destroy();
    }
});