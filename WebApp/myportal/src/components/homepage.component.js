import React, { Component } from 'react';
import Course from './course.component'
import ConfirmedCourse from './confirmCourse.component'
import EnrolledCourse from './enrolledCourse.components'
// import SuccessCourse from './successCourse.component'
import "../App.css"
const axios = require('axios');

export default class HomePage extends Component {
    constructor(props) {
        super(props); 

        this.onChangeCourseId = this.onChangeCourseId.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
        this.selectCourse = this.selectCourse.bind(this);
        this.cancelConfirm = this.cancelConfirm.bind(this)
        this.confirmConfirm = this.confirmConfirm.bind(this);

        this.state = {
            alertType: "",
            alert: false,
            subject: '', 
            courseId: '', 
            career: '', 
            addedToCartCourses: [],
            enrolledCourses: [],
            enrolledCoursesString: '',
            successCourses: [],
            searchStage: true,
            selectStage: false, 
            confirmStage: false, 
            notification: "",
            courses: [
                {
                    _id: '1',
                    class: "1057",
                    section: "CH01-CLB Regular",
                    dayTime: "Tu 15:00 - 17:00", 
                    room: "Think Tank 13 (1.508)", 
                    instructor: "Staff", 
                    meetDate: "20/05/2019 - 16/08/2019",
                    status: 'Available', 
                    units: '24.00',
                    seats:  0
                },
                // {   
                //     _id: '2',
                //     class: "1194",
                //     section: "CH02-CLB Regular",
                //     dayTime: "Th 15:00 - 17:00", 
                //     room: "Think Tank 13 (1.508)", 
                //     instructor: "Staff", 
                //     meetDate: "20/05/2019 - 16/08/2019",
                //     status: 'Available', 
                //     units: '24.00' 
                // }
            ],
        }
        
    }

    onChangeSubject(e) {
        this.setState({
            subject: e.target.value
        }) 
     }
 
     onChangeCourseId(e) {
         this.setState({
            courseId: e.target.value
         }) 
      }
 
      onChangeCourseCareer(e) {
         this.setState({
            career: e.target.value
         }) 
     }

     cancelConfirm() {
         let tempEnrolledCourses = [...this.state.enrolledCourses]
         if (this.state.enrolledCourses.length == 0) {
            tempEnrolledCourses = []
         }
         this.setState({
            searchStage: true,
            selectStage: false, 
            confirmStage: false,
            successStage: false,
            enrolledCourses: tempEnrolledCourses,
            notification: "" 
         })
     }

    confirmConfirm() {
        let success = true 
        let tempEnrolledCourses
        if (success) {
            tempEnrolledCourses = [...this.state.enrolledCourses]
            tempEnrolledCourses.push(this.state.courses[0])
        }
        this.setState({
           searchStage: true,
           selectStage: false, 
           confirmStage: false,
           alert: true,
           alertType: "confirm",
           addedToCartCourses: [],
           enrolledCourses: tempEnrolledCourses,
           notification: `You are successfully enrolled into Course ${this.state.courses[0].class}`
        })
        this.timer = setInterval(() => {this.setState({success: false})}, 3000)
    }

    selectCourse(course) {
        const tempAddedToCart = [...this.state.addedToCartCourses];
        tempAddedToCart.push(course)
        this.setState({
            addedToCartCourses: tempAddedToCart,
            selectStage: false, 
            alert: true, 
            alertType: "addToCart",
            confirmStage: true , 
            notification: `Course ${course.class} is added to cart` 
        })
        this.timer = setInterval(() => {this.setState({addToCart: false})}, 3000)
    }

    // deleteCourse(course) {
    //     this.setState({
    //         enrolledCourses: this.state.enrolledCourses.filter(c => c._id !== course._id),
    //         notification: `${course.class} is removed from cart` 
    //     })
    // }

    selectCoursesList() {
        return this.state.courses.map(course => {
            return <Course course={course} selectCourse={this.selectCourse} key={course._id}/>;
        })
    }

    enrolledCoursesList() {
        return this.state.enrolledCourses.map(course => {
            return <EnrolledCourse course={course} key={course._id}/>;
        })
    }

    confirmedCourseList() {
        return this.state.addedToCartCourses.map(course => {
            return <ConfirmedCourse course={course}  key={course._id}/>;
        })
    }


    onSubmit(e) { 
        e.preventDefault()
        let course = [...this.state.courses]
        axios.get('http://localhost:3001/read-from-node?courseid=' + this.state.courseId)
        .then(response => {
            console.log(response)
            course[0].seats = response.data
            this.setState({
                selectStage: true, 
                searchStage: false,
                confirmStage: false, 
                courses: course,
                courseId: "",
                alert: false,
            })
        })
        .catch(error => {
            this.setState({
                alert: true,
                alertType: "invalidID",
                notification: "Invalid Course Id"
            });
        })
        this.timer = setInterval(() => {this.setState({alert: false})}, 3000)
    }

    render() {
        let alertMessage
        let enrolmentsummary
        let alertSignal = ""
        if (this.state.alert) {
            if (this.state.alertType == "invalidID"){
              alertSignal = "alert alert-danger"
            } else if (this.state.alertType == "addToCart") {
                alertSignal = "alert alert-primary"
            } else {
                alertSignal = "alert alert-success"
            }
            alertMessage = <div class={alertSignal} role="alert">
                            {this.state.notification}
            </div> 
        } 
        
        if (this.state.enrolledCourses.length != 0)  {
            enrolmentsummary = <div>
            <h3>Enrollment Summary</h3>
            <div></div>
            <table className="table">
            <thead className="thead-light">
                <tr>
                <th>Class</th>
                <th>Days/Times</th>
                <th>Room</th>
                <th>Instructor</th>
                <th>Units</th>
                <th>Status</th>
                </tr>
            </thead>
            <tbody>
                {this.enrolledCoursesList()}
            </tbody>
            </table>
            </div>
        } else {
            enrolmentsummary = <div className="summary">
                    <h2>No Enrolled Courses</h2>
                </div>
        }

        const notification = this.state.notification
        if (this.state.searchStage) {
            let enrolledClass = <b></b>; 
            if (this.state.enrolledCoursesString != "") {
                enrolledClass = <b>Enrolled class:</b>
            }
            return (
                <div>
                    {alertMessage}
                    {/* {successMessage} */}
                    <div className = 'row'>
                        <div className ='col'>
                           {enrolledClass} {this.state.enrolledCoursesString}
                        </div>
                    </div>
                    <div className = 'row'>
                    <div className="col-4">
                    <h2>Class Search</h2> 
                    <form onSubmit={this.onSubmit}>
                        <div className = "form-group"> 
                            <input type="text" class="form-control" placeholder="Enter Course ID"  value={this.state.courseId}
                        onChange={this.onChangeCourseId}/>        
                        </div>
                            <div className="form-group">
                                <input type="submit" value="Search" className="btn btn-primary"/>
                            </div>
                    </form>
                    </div>
                    <div className="col-8">
                        {enrolmentsummary}
                    </div>
                    </div> 
                </div>
            )
        }
        else if (this.state.selectStage) {
            return (
                <div>
                    {alertMessage}
                     <h3>Available Courses</h3>
                     <div></div>
                     <table className="table">
                     <thead className="thead-light">
                         <tr>
                         <th>Class</th>
                         <th>Section</th>
                         <th>Days & Times</th>
                         <th>Room</th>
                         <th>Instructor</th>
                         <th>Meeting Dates</th>
                         <th>Status</th>
                         <th>Seats</th>
                         <th></th>
                         </tr>
                     </thead>
                     <tbody>
                         {this.selectCoursesList()}
                        
                     </tbody>
                     </table>
                     <div className="div-right">
                        <button type="button" className="btn btn-danger cancel-right"  onClick={this.cancelConfirm}>
                                    Back to Menu
                        </button>
                     </div>
                </div>
             )
        } 
        else if (this.state.confirmStage) {
            return (
                <div>
                    {alertMessage}
                     <h3>Confirm Classes</h3>
                     <div></div>
                     <table className="table">
                        <thead className="thead-light">
                            <tr>
                            <th>Class</th>
                            <th>Days/Times</th>
                            <th>Room</th>
                            <th>Instructor</th>
                            <th>Units</th>
                            <th>Status</th>
                            <th>Seats</th>
                            </tr>
                        </thead>
                        <tbody>
                            {this.confirmedCourseList()}
                        </tbody>
                        </table>
                        <div className="d-flex flex-row-reverse">
                            <div className="p-2">
                            <button type="button" className="btn btn-danger confirm-flex" onClick={this.cancelConfirm}>
                               Back to Menu
                            </button>
                            <button type="button" className="btn btn-primary confirm-flex" onClick={this.confirmConfirm}>
                                    Finished Enrolling
                            </button>
                            </div>
                        </div>
                </div>
             )
        }
    }
}


