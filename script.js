// Discord Server Cloner JavaScript
class DiscordCloner {
    constructor() {
        this.isCloning = false;
        this.progress = 0;
        this.stats = {
            channels: 0,
            roles: 0,
            emojis: 0,
            stickers: 0
        };
        
        this.initializeEventListeners();
    }

    initializeEventListeners() {
        const cloneBtn = document.getElementById('clone-btn');
        const tokenInput = document.getElementById('discord-token');
        const sourceInput = document.getElementById('source-server');
        const targetInput = document.getElementById('target-server');

        cloneBtn.addEventListener('click', () => this.startCloning());
        
        // Real-time validation
        [tokenInput, sourceInput, targetInput].forEach(input => {
            input.addEventListener('input', () => this.validateInputs());
        });
        
        this.validateInputs();
    }

    validateInputs() {
        const token = document.getElementById('discord-token').value;
        const sourceId = document.getElementById('source-server').value;
        const targetId = document.getElementById('target-server').value;
        const cloneBtn = document.getElementById('clone-btn');

        const isValid = token.length > 50 && 
                       sourceId.length >= 17 && 
                       targetId.length >= 17 &&
                       sourceId !== targetId;

        cloneBtn.disabled = !isValid || this.isCloning;
    }

    async startCloning() {
        if (this.isCloning) return;
        
        this.isCloning = true;
        this.progress = 0;
        this.resetStats();
        
        const cloneBtn = document.getElementById('clone-btn');
        cloneBtn.disabled = true;
        cloneBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> CLONING...';
        
        this.logMessage('🚀 Initializing Discord API connection...');
        await this.delay(1000);
        
        this.logMessage('✅ Connection established successfully!');
        await this.delay(500);
        
        this.logMessage('📡 Fetching source server information...');
        await this.delay(1500);
        
        this.logMessage('🔍 Analyzing server structure...');
        await this.delay(1000);
        
        // Simulate cloning process
        await this.simulateCloning();
        
        this.logMessage('🎉 Server cloning completed successfully!');
        this.logMessage('✨ All data has been transferred to the target server.');
        
        cloneBtn.innerHTML = '<i class="fas fa-check"></i> COMPLETED';
        cloneBtn.style.background = 'linear-gradient(90deg, #22c55e, #16a34a)';
        
        setTimeout(() => {
            this.resetCloner();
        }, 3000);
    }

    async simulateCloning() {
        const steps = [
            { name: 'Creating channels', type: 'channels', count: 15, icon: '📝' },
            { name: 'Copying roles', type: 'roles', count: 8, icon: '👥' },
            { name: 'Transferring emojis', type: 'emojis', count: 25, icon: '😀' },
            { name: 'Moving stickers', type: 'stickers', count: 12, icon: '🎨' }
        ];

        for (let i = 0; i < steps.length; i++) {
            const step = steps[i];
            this.logMessage(`${step.icon} ${step.name}...`);
            
            // Simulate gradual progress for each step
            for (let j = 0; j <= step.count; j++) {
                this.stats[step.type] = j;
                this.progress = ((i * 25) + (j / step.count * 25));
                this.updateUI();
                await this.delay(100);
            }
            
            this.logMessage(`✅ ${step.name} completed (${step.count} items)`);
            await this.delay(300);
        }
    }

    updateUI() {
        // Update statistics
        document.getElementById('channels-count').textContent = this.stats.channels;
        document.getElementById('roles-count').textContent = this.stats.roles;
        document.getElementById('emojis-count').textContent = this.stats.emojis;
        document.getElementById('stickers-count').textContent = this.stats.stickers;
        
        // Update progress bar
        const progressFill = document.getElementById('progress-fill');
        const progressPercent = document.getElementById('progress-percent');
        
        progressFill.style.width = `${this.progress}%`;
        progressPercent.textContent = `${Math.round(this.progress)}%`;
    }

    logMessage(message) {
        const monitorContent = document.getElementById('monitor-content');
        const logEntry = document.createElement('div');
        logEntry.className = 'monitor-log';
        logEntry.innerHTML = `<span style="color: #666;">[${new Date().toLocaleTimeString()}]</span> ${message}`;
        
        monitorContent.appendChild(logEntry);
        monitorContent.scrollTop = monitorContent.scrollHeight;
        
        // Keep only last 10 messages
        const logs = monitorContent.querySelectorAll('.monitor-log');
        if (logs.length > 10) {
            logs[0].remove();
        }
    }

    resetStats() {
        this.stats = { channels: 0, roles: 0, emojis: 0, stickers: 0 };
        this.updateUI();
    }

    resetCloner() {
        this.isCloning = false;
        this.progress = 0;
        
        const cloneBtn = document.getElementById('clone-btn');
        cloneBtn.innerHTML = '<i class="fas fa-copy"></i> CLONE SERVER';
        cloneBtn.style.background = 'linear-gradient(90deg, #5865f2, #7289da)';
        
        this.validateInputs();
        
        const monitorContent = document.getElementById('monitor-content');
        monitorContent.innerHTML = `
            <div class="monitor-status">
                <i class="fas fa-clock"></i>
                <span>READY TO CLONE YOUR SERVER...</span>
            </div>
        `;
        
        this.resetStats();
    }

    delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
}

// Initialize the cloner when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new DiscordCloner();
    
    // Add some visual effects
    addVisualEffects();
});

function addVisualEffects() {
    // Animated background particles
    createParticles();
    
    // Smooth scrolling for internal links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({ behavior: 'smooth' });
            }
        });
    });
    
    // Add hover effects to stat cards
    document.querySelectorAll('.stat-card').forEach(card => {
        card.addEventListener('mouseenter', function() {
            this.style.transform = 'translateY(-5px) scale(1.02)';
        });
        
        card.addEventListener('mouseleave', function() {
            this.style.transform = 'translateY(0) scale(1)';
        });
    });
}

function createParticles() {
    const particleContainer = document.createElement('div');
    particleContainer.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        pointer-events: none;
        z-index: -1;
    `;
    
    document.body.appendChild(particleContainer);
    
    for (let i = 0; i < 50; i++) {
        createParticle(particleContainer);
    }
}

function createParticle(container) {
    const particle = document.createElement('div');
    particle.style.cssText = `
        position: absolute;
        width: 2px;
        height: 2px;
        background: rgba(88, 101, 242, 0.3);
        border-radius: 50%;
        animation: float ${5 + Math.random() * 10}s infinite linear;
    `;
    
    particle.style.left = Math.random() * 100 + '%';
    particle.style.animationDelay = Math.random() * 10 + 's';
    
    container.appendChild(particle);
    
    // Add floating animation
    const style = document.createElement('style');
    style.textContent = `
        @keyframes float {
            0% {
                transform: translateY(100vh) rotate(0deg);
                opacity: 0;
            }
            10% {
                opacity: 1;
            }
            90% {
                opacity: 1;
            }
            100% {
                transform: translateY(-100px) rotate(360deg);
                opacity: 0;
            }
        }
    `;
    
    if (!document.querySelector('#particle-styles')) {
        style.id = 'particle-styles';
        document.head.appendChild(style);
    }
}

// Add some Easter eggs
document.addEventListener('keydown', (e) => {
    // Konami code easter egg
    const konamiCode = [38, 38, 40, 40, 37, 39, 37, 39, 66, 65];
    if (!window.konamiIndex) window.konamiIndex = 0;
    
    if (e.keyCode === konamiCode[window.konamiIndex]) {
        window.konamiIndex++;
        if (window.konamiIndex === konamiCode.length) {
            triggerEasterEgg();
            window.konamiIndex = 0;
        }
    } else {
        window.konamiIndex = 0;
    }
});

function triggerEasterEgg() {
    const body = document.body;
    body.style.animation = 'rainbow 2s infinite';
    
    const style = document.createElement('style');
    style.textContent = `
        @keyframes rainbow {
            0% { filter: hue-rotate(0deg); }
            100% { filter: hue-rotate(360deg); }
        }
    `;
    document.head.appendChild(style);
    
    setTimeout(() => {
        body.style.animation = '';
        style.remove();
    }, 5000);
    
    // Show secret message
    alert('🎉 Easter egg activated! You found the secret code!');
}