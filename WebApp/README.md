# Subject Enrolment System

Web Application for SUTD Subject Enrolment System, created using React.

## How to Setup
* Run the command 'npm run start'

## How to Use

### Add to Enrollment Cart
The home page is as shown below. Use the search bar to search for the desired CourseID.
<p align="center"><img src="../pictures/home-page.png" /></p>
If the courseID matches an existing course, the course description with its number of available seats will be displayed.
<p align="center"><img src="../pictures/select-page.png" /></p>
Click 'Select' to add to cart.

### Finish Enrolling
To finish enrolling for the selected course(s), click on the 'Finished Enrolling' button.
<p align="center"><img src="../pictures/finish-enrolling-page.png" /></p>
The database will indicate whether the enrollment is successful:
<p align="center"><img src="../pictures/successful-enrollment.png" /></p>
...or not:
<p align="center"><img src="../pictures/connection-refused.png" /></p>
