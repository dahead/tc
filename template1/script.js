// Global variables
let isModalOpen = false;

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    initializeTagCloud();
    setupEventListeners();
});

function initializeTagCloud() {
    // Animate tags on load
    const tags = document.querySelectorAll('.tag');
    tags.forEach(function(tag, index) {
        setTimeout(function() {
            tag.style.transition = 'all 0.5s ease';
            tag.style.opacity = '1';
            tag.style.transform = 'scale(1)';
            tag.classList.add('loaded');
        }, index * 50);

        // Add click event listener
        tag.addEventListener('click', function() {
            const word = this.getAttribute('data-word');
            const filesJson = this.getAttribute('data-files');
            const files = JSON.parse(filesJson);
            showFiles(word, files);
        });
    });

    console.log('Tag cloud initialized with', tags.length, 'tags');
}

function setupEventListeners() {
    // Close modal when clicking outside
    window.onclick = function(event) {
        const modal = document.getElementById('fileModal');
        const videoModal = document.getElementById('videoModal');
        if (event.target === modal) {
            closeModal();
        }
        if (event.target === videoModal) {
            closeVideoModal();
        }
    };

    // Keyboard shortcuts
    document.addEventListener('keydown', function(event) {
        if (event.key === 'Escape' && isModalOpen) {
            closeModal();
            closeVideoModal();
        }
    });
}

function showFiles(word, files) {
    console.log('Showing files for word:', word, files);
    document.getElementById('modalTitle').textContent = 'Files containing "' + word + '"';

    const fileListContainer = document.getElementById('fileListContainer');
    fileListContainer.innerHTML = '';

    // Create vertical list of file tags
    files.forEach(function(file, index) {
        const fileTag = document.createElement('div');
        fileTag.className = 'file-tag';

        // Extract just the filename without path
        const filename = file.split('/').pop() || file.split('\\').pop() || file;
        fileTag.textContent = filename;
        fileTag.title = file; // Show full path on hover

        // Add click event to copy to clipboard or play video
        fileTag.addEventListener('click', function() {
            if (isVideoFile(file)) {
                playVideo(file);
            } else {
                copyToClipboard(file, fileTag);
            }
        });

        fileListContainer.appendChild(fileTag);

        // Animate file tag appearance
        setTimeout(function() {
            fileTag.classList.add('loaded');
        }, index * 100);
    });

    document.getElementById('fileModal').style.display = 'block';
    isModalOpen = true;
}

function isVideoFile(filename) {
    const videoExtensions = ['.mp4', '.webm', '.ogg', '.avi', '.mov', '.wmv', '.flv', '.mkv'];
    return videoExtensions.some(ext => filename.toLowerCase().endsWith(ext));
}

// open video in new tab
function playVideo(filepath) {
    window.open(filepath, '_blank');
}

// open video in small player
// function playVideo(filepath) {
//     const filename = filepath.split('/').pop() || filepath.split('\\').pop() || filepath;
//     document.getElementById('videoModalTitle').textContent = filename;
//     document.getElementById('videoPlayer').src = filepath;
//     document.getElementById('videoModal').style.display = 'block';
// }

function closeVideoModal() {
    const video = document.getElementById('videoPlayer');
    video.pause();
    video.src = '';
    document.getElementById('videoModal').style.display = 'none';
}

async function copyToClipboard(text, element) {
    try {
        await navigator.clipboard.writeText(text);

        // Visual feedback
        element.classList.add('copied');
        setTimeout(() => element.classList.remove('copied'), 1000);

        showNotification('Copied: ' + (text.split('/').pop() || text.split('\\').pop() || text));
    } catch (err) {
        console.error('Failed to copy text: ', err);

        // Fallback for older browsers
        const textArea = document.createElement('textarea');
        textArea.value = text;
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();

        try {
            document.execCommand('copy');
            showNotification('Copied: ' + (text.split('/').pop() || text.split('\\').pop() || text));
            element.classList.add('copied');
            setTimeout(() => element.classList.remove('copied'), 1000);
        } catch (err) {
            showNotification('Failed to copy to clipboard');
        }

        document.body.removeChild(textArea);
    }
}

function closeModal() {
    document.getElementById('fileModal').style.display = 'none';
    isModalOpen = false;
}

function showNotification(message) {
    // Remove existing notification if any
    const existing = document.querySelector('.notification');
    if (existing) {
        existing.remove();
    }

    // Create notification element
    const notification = document.createElement('div');
    notification.className = 'notification';
    notification.textContent = message;

    document.body.appendChild(notification);

    // Remove after 3 seconds
    setTimeout(function() {
        if (document.body.contains(notification)) {
            document.body.removeChild(notification);
        }
    }, 3000);
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Utility functions for potential future features
function filterTags(searchTerm) {
    const tags = document.querySelectorAll('.tag');
    tags.forEach(function(tag) {
        const text = tag.textContent.toLowerCase();
        if (text.includes(searchTerm.toLowerCase())) {
            tag.style.display = 'block';
        } else {
            tag.style.display = 'none';
        }
    });
}

function resetTagFilter() {
    const tags = document.querySelectorAll('.tag');
    tags.forEach(function(tag) {
        tag.style.display = 'block';
    });
}