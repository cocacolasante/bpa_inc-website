// ===================================
// Navigation Functionality
// ===================================
document.addEventListener('DOMContentLoaded', function() {
    const navbar = document.getElementById('navbar');
    const mobileMenuToggle = document.getElementById('mobileMenuToggle');
    const navMenu = document.getElementById('navMenu');

    // Navbar scroll effect
    let lastScroll = 0;
    window.addEventListener('scroll', function() {
        const currentScroll = window.pageYOffset;

        if (currentScroll > 100) {
            navbar.classList.add('scrolled');
        } else {
            navbar.classList.remove('scrolled');
        }

        lastScroll = currentScroll;
    });

    // Mobile menu toggle
    if (mobileMenuToggle) {
        mobileMenuToggle.addEventListener('click', function() {
            this.classList.toggle('active');
            navMenu.classList.toggle('active');
        });

        // Close mobile menu when clicking on a link
        const navLinks = navMenu.querySelectorAll('.nav-link');
        navLinks.forEach(link => {
            link.addEventListener('click', function() {
                mobileMenuToggle.classList.remove('active');
                navMenu.classList.remove('active');
            });
        });

        // Close mobile menu when clicking outside
        document.addEventListener('click', function(event) {
            const isClickInsideNav = navMenu.contains(event.target);
            const isClickOnToggle = mobileMenuToggle.contains(event.target);

            if (!isClickInsideNav && !isClickOnToggle && navMenu.classList.contains('active')) {
                mobileMenuToggle.classList.remove('active');
                navMenu.classList.remove('active');
            }
        });
    }
});

// ===================================
// Smooth Scrolling
// ===================================
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function(e) {
        const href = this.getAttribute('href');
        if (href !== '#' && href.length > 1) {
            e.preventDefault();
            const target = document.querySelector(href);
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        }
    });
});

// ===================================
// Contact Form Handling
// ===================================
const contactForm = document.getElementById('contactForm');

if (contactForm) {
    contactForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        // Clear previous errors
        clearErrors();

        // Get form data
        const formData = {
            full_name: document.getElementById('fullName').value.trim(),
            business_name: document.getElementById('businessName').value.trim(),
            email: document.getElementById('email').value.trim(),
            phone: document.getElementById('phone').value.trim(),
            website: document.getElementById('website').value.trim(),
            message: document.getElementById('message').value.trim()
        };

        // Client-side validation
        if (!validateForm(formData)) {
            return;
        }

        // Show loading state
        const submitButton = document.getElementById('submitButton');
        const buttonText = submitButton.querySelector('.button-text');
        const buttonLoading = submitButton.querySelector('.button-loading');

        submitButton.disabled = true;
        buttonText.style.display = 'none';
        buttonLoading.style.display = 'flex';

        try {
            const response = await fetch('/api/contact', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            });

            const data = await response.json();

            if (response.ok && data.success) {
                showMessage(data.message, 'success');
                contactForm.reset();
            } else {
                showMessage(data.message || 'An error occurred. Please try again.', 'error');
            }
        } catch (error) {
            console.error('Form submission error:', error);
            showMessage('An error occurred. Please try again later.', 'error');
        } finally {
            // Reset button state
            submitButton.disabled = false;
            buttonText.style.display = 'block';
            buttonLoading.style.display = 'none';
        }
    });
}

function validateForm(data) {
    let isValid = true;

    // Full Name validation
    if (!data.full_name || data.full_name.length < 2) {
        showError('fullNameError', 'Please enter your full name (at least 2 characters)');
        isValid = false;
    }

    // Business Name validation
    if (!data.business_name) {
        showError('businessNameError', 'Please enter your business name');
        isValid = false;
    }

    // Email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!data.email || !emailRegex.test(data.email)) {
        showError('emailError', 'Please enter a valid email address');
        isValid = false;
    }

    // Phone validation
    const phoneRegex = /^[\d\s\-\+\(\)]{10,}$/;
    if (!data.phone || !phoneRegex.test(data.phone)) {
        showError('phoneError', 'Please enter a valid phone number');
        isValid = false;
    }

    // Message validation
    if (!data.message || data.message.length < 10) {
        showError('messageError', 'Please enter a message (at least 10 characters)');
        isValid = false;
    }

    return isValid;
}

function showError(elementId, message) {
    const errorElement = document.getElementById(elementId);
    if (errorElement) {
        errorElement.textContent = message;
        errorElement.style.display = 'block';
    }
}

function clearErrors() {
    const errorElements = document.querySelectorAll('.error-message');
    errorElements.forEach(element => {
        element.textContent = '';
        element.style.display = 'none';
    });
}

function showMessage(message, type) {
    const formMessage = document.getElementById('formMessage');
    if (formMessage) {
        formMessage.textContent = message;
        formMessage.className = 'form-message ' + type;
        formMessage.style.display = 'block';

        // Hide message after 5 seconds for success, keep error messages visible
        if (type === 'success') {
            setTimeout(() => {
                formMessage.style.display = 'none';
            }, 5000);
        }
    }
}

// ===================================
// Intersection Observer for Animations
// ===================================
const observerOptions = {
    threshold: 0.1,
    rootMargin: '0px 0px -50px 0px'
};

const observer = new IntersectionObserver(function(entries) {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.style.opacity = '1';
            entry.target.style.transform = 'translateY(0)';
        }
    });
}, observerOptions);

// Observe elements that should animate on scroll
document.addEventListener('DOMContentLoaded', function() {
    const animateElements = document.querySelectorAll('.feature-card, .service-card, .portfolio-card, .founder-card, .value-card, .stat-item');

    animateElements.forEach(element => {
        element.style.opacity = '0';
        element.style.transform = 'translateY(20px)';
        element.style.transition = 'opacity 0.6s ease, transform 0.6s ease';
        observer.observe(element);
    });
});

// ===================================
// Number Counter Animation (for stats)
// ===================================
function animateCounter(element, target, duration = 2000) {
    const start = 0;
    const increment = target / (duration / 16);
    let current = start;

    const timer = setInterval(() => {
        current += increment;
        if (current >= target) {
            element.textContent = target;
            clearInterval(timer);
        } else {
            element.textContent = Math.floor(current);
        }
    }, 16);
}

// Animate stat numbers when they come into view
const statsObserver = new IntersectionObserver(function(entries) {
    entries.forEach(entry => {
        if (entry.isIntersecting && !entry.target.classList.contains('animated')) {
            const statNumber = entry.target.querySelector('.stat-number');
            if (statNumber) {
                const value = statNumber.textContent;
                const numericValue = parseInt(value.replace(/\D/g, ''));

                if (!isNaN(numericValue)) {
                    statNumber.textContent = '0';
                    animateCounter(statNumber, numericValue);
                }

                entry.target.classList.add('animated');
            }
        }
    });
}, { threshold: 0.5 });

document.addEventListener('DOMContentLoaded', function() {
    const statItems = document.querySelectorAll('.stat-item');
    statItems.forEach(item => statsObserver.observe(item));
});

// ===================================
// Form Input Focus Effects
// ===================================
document.addEventListener('DOMContentLoaded', function() {
    const formInputs = document.querySelectorAll('.contact-form input, .contact-form textarea');

    formInputs.forEach(input => {
        input.addEventListener('focus', function() {
            this.parentElement.classList.add('focused');
        });

        input.addEventListener('blur', function() {
            if (!this.value) {
                this.parentElement.classList.remove('focused');
            }
        });

        // Clear error on input
        input.addEventListener('input', function() {
            const errorElement = this.parentElement.querySelector('.error-message');
            if (errorElement) {
                errorElement.style.display = 'none';
            }
        });
    });
});

// ===================================
// Scroll Progress Indicator (optional enhancement)
// ===================================
window.addEventListener('scroll', function() {
    const winScroll = document.body.scrollTop || document.documentElement.scrollTop;
    const height = document.documentElement.scrollHeight - document.documentElement.clientHeight;
    const scrolled = (winScroll / height) * 100;

    // You can add a progress bar element if desired
    // document.getElementById('progressBar').style.width = scrolled + '%';
});

// ===================================
// Lazy Loading Images (if images are added later)
// ===================================
if ('IntersectionObserver' in window) {
    const imageObserver = new IntersectionObserver(function(entries, observer) {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const img = entry.target;
                if (img.dataset.src) {
                    img.src = img.dataset.src;
                    img.classList.remove('lazy');
                    imageObserver.unobserve(img);
                }
            }
        });
    });

    document.querySelectorAll('img.lazy').forEach(img => {
        imageObserver.observe(img);
    });
}

// ===================================
// Utility: Debounce Function
// ===================================
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// ===================================
// Performance: Reduce scroll event frequency
// ===================================
const debouncedScroll = debounce(function() {
    // Any scroll-based animations or checks can go here
}, 100);

window.addEventListener('scroll', debouncedScroll);

// ===================================
// Pricing Toggle — Monthly / Annual
// ===================================
const billingToggle = document.getElementById('billingToggle');

if (billingToggle) {
    billingToggle.addEventListener('change', function() {
        const isAnnual = this.checked;
        const amounts = document.querySelectorAll('.monthly-amount');

        amounts.forEach(function(el) {
            const monthly = parseInt(el.dataset.monthly, 10);
            const annual = parseInt(el.dataset.annual, 10);
            const display = isAnnual ? annual : monthly;
            el.textContent = '$' + display.toLocaleString();
        });

        const badge = document.querySelector('.toggle-badge');
        if (badge) {
            badge.style.opacity = isAnnual ? '1' : '0.5';
        }
    });
}

// ===================================
// Audit Form Handling
// ===================================
const auditForm = document.getElementById('auditForm');

if (auditForm) {
    auditForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const submitBtn = document.getElementById('auditSubmitButton');
        const btnText = submitBtn.querySelector('.button-text');
        const btnLoading = submitBtn.querySelector('.button-loading');
        submitBtn.disabled = true;
        btnText.style.display = 'none';
        btnLoading.style.display = 'flex';

        const formData = {
            full_name:     document.getElementById('auditFullName').value.trim(),
            business_name: document.getElementById('auditBusinessName').value.trim(),
            email:         document.getElementById('auditEmail').value.trim(),
            phone:         document.getElementById('auditPhone').value.trim(),
            website:       '',
            message: [
                'Team size: ' + (document.getElementById('auditTeamSize').value || 'Not specified'),
                'Hours/week on manual tasks: ' + (document.getElementById('auditHoursWeek').value || 'Not specified'),
                'Pain points: ' + document.getElementById('auditPainPoints').value.trim(),
            ].join('\n'),
            source: 'roi_audit',
        };

        try {
            const response = await fetch('/api/contact', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(formData),
            });
            const data = await response.json();
            const msg = document.getElementById('auditFormMessage');
            if (response.ok && data.success) {
                msg.textContent = "Thanks! We'll have your audit ready within 24 hours. Check your email.";
                msg.className = 'form-message success';
                msg.style.display = 'block';
                auditForm.reset();
            } else {
                msg.textContent = data.message || 'An error occurred. Please try again.';
                msg.className = 'form-message error';
                msg.style.display = 'block';
            }
        } catch (err) {
            const msg = document.getElementById('auditFormMessage');
            msg.textContent = 'An error occurred. Please try again.';
            msg.className = 'form-message error';
            msg.style.display = 'block';
        } finally {
            submitBtn.disabled = false;
            btnText.style.display = 'block';
            btnLoading.style.display = 'none';
        }
    });
}

// ===================================
// Partner Application Form
// ===================================
const partnerForm = document.getElementById('partnerForm');

if (partnerForm) {
    partnerForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const submitBtn = document.getElementById('partnerSubmitButton');
        const btnText = submitBtn.querySelector('.button-text');
        const btnLoading = submitBtn.querySelector('.button-loading');
        submitBtn.disabled = true;
        btnText.style.display = 'none';
        btnLoading.style.display = 'flex';

        const formData = {
            company_name:      document.getElementById('partnerCompanyName').value.trim(),
            contact_name:      document.getElementById('partnerContactName').value.trim(),
            contact_email:     document.getElementById('partnerEmail').value.trim(),
            contact_phone:     document.getElementById('partnerPhone').value.trim(),
            website:           document.getElementById('partnerWebsite').value.trim(),
            years_in_business: document.getElementById('partnerYears').value,
            client_count:      document.getElementById('partnerClients').value,
            expected_volume:   document.getElementById('partnerVolume').value,
            why_partner:       document.getElementById('partnerWhy').value.trim(),
        };

        try {
            const response = await fetch('/api/partner/apply', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(formData),
            });
            const data = await response.json();
            const msg = document.getElementById('partnerFormMessage');
            msg.textContent = data.message;
            msg.className = 'form-message ' + (data.success ? 'success' : 'error');
            msg.style.display = 'block';
            if (data.success) {
                partnerForm.reset();
                window.scrollTo({ top: msg.offsetTop - 100, behavior: 'smooth' });
            }
        } catch (err) {
            const msg = document.getElementById('partnerFormMessage');
            msg.textContent = 'An error occurred. Please try again.';
            msg.className = 'form-message error';
            msg.style.display = 'block';
        } finally {
            submitBtn.disabled = false;
            btnText.style.display = 'block';
            btnLoading.style.display = 'none';
        }
    });
}
