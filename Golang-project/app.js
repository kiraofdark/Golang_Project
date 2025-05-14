// Fetch golf courses data from the API
fetch('http://localhost:5000/api/golfcourses')
    .then(response => response.json())
    .then(data => {
        // Create a list of courses and display it
        const coursesList = document.getElementById('courses-list');
        if (data.length === 0) {
            coursesList.innerHTML = '<p>No courses available.</p>';
        } else {
            let htmlContent = '<ul>';
            data.forEach(course => {
                htmlContent += `
                    <li>
                        <strong>${course.coursename}</strong><br>
                        Price: ${course.price} Bath<br>
                        Total Holes: ${course.totalhole}<br>
                    </li>
                `;
            });
            htmlContent += '</ul>';
            coursesList.innerHTML = htmlContent;
        }
    })
    .catch(error => {
        console.error('Error fetching data:', error);
        document.getElementById('courses-list').innerHTML = '<p>Error loading courses.</p>';
    });
