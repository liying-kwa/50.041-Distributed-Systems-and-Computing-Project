import React, { Component } from 'react';
import Course from './course.component'
import ConfirmedCourse from './confirmCourse.component'
import EnrolledCourse from './enrolledCourse.components'
import SuccessCourse from './successCourse.component'
import "../App.css"

export default class HomePage extends Component {
    constructor(props) {
        super(props); 

        this.onChangeCouseNumber = this.onChangeCouseNumber.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
        this.selectCourse = this.selectCourse.bind(this);
        this.deleteCourse = this.deleteCourse.bind(this);
        this.goToStep2 = this.goToStep2.bind(this);
        this.cancelConfirm = this.cancelConfirm.bind(this)
        this.previousConfirm = this.previousConfirm.bind(this);
        this.confirmConfirm = this.confirmConfirm.bind(this);
        this.addAnotherClass = this.addAnotherClass.bind(this);

        this.state = {
            subject: '', 
            number: '', 
            career: '', 
            enrolledCourses: [],
            enrolledCoursesString: '',
            successCourses: [],
            searchStage: true,
            selectStage: false, 
            confirmStage: false, 
            successStage: false,
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
                    units: '24.00'
                },
                {   
                    _id: '2',
                    class: "1194",
                    section: "CH02-CLB Regular",
                    dayTime: "Th 15:00 - 17:00", 
                    room: "Think Tank 13 (1.508)", 
                    instructor: "Staff", 
                    meetDate: "20/05/2019 - 16/08/2019",
                    status: 'Available', 
                    units: '24.00'
                }
            ],
            // search: false,
        }
        
    }

    onChangeSubject(e) {
        this.setState({
            subject: e.target.value
        }) 
     }
 
     onChangeCouseNumber(e) {
         this.setState({
            number: e.target.value
         }) 
      }
 
      onChangeCourseCareer(e) {
         this.setState({
            career: e.target.value
         }) 
     }

     goToStep2() {
         this.setState({
            searchStage: false,
            selectStage: false, 
            confirmStage: true, 
            successStage: false,
            notification: ""
         })
     }

     cancelConfirm() {
         this.setState({
            searchStage: true,
            selectStage: false, 
            confirmStage: false,
            successStage: false,
            enrolledCourses: [],
            notification: "" 
         })
     }

     previousConfirm() {
        this.setState({
           searchStage: true,
           selectStage: false, 
           confirmStage: false,
           successStage: false,
           notification: ""
        })
    }

    confirmConfirm() {
        this.setState({
           searchStage: false,
           selectStage: false, 
           confirmStage: false,
           successStage: true,
           successCourses: [...this.state.enrolledCourses],
           notification: ""
        })
        let result = this.state.enrolledCourses.map(x => x.class)
        let array = result + "," + this.state.enrolledCoursesString
        console.log(array)
        const requestOptions = {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({"value": array })
        };
        fetch('http://localhost:3000/api/v1/student/1003123', requestOptions)
            .then(response => response.json())
            .then(data => console.log(data));
    }

    addAnotherClass() {
        this.setState({
            searchStage: true,
            selectStage: false, 
            confirmStage: false,
            successStage: false,
            enrolledCourses: [],
            successCourses: [...this.state.enrolledCourses],
            notification: ""
         })
    }

    selectCourse(course) {
        const enrolledCourses = [...this.state.enrolledCourses];
        enrolledCourses.push(course)
        this.setState({
            enrolledCourses: enrolledCourses,
            selectStage: false, 
            searchStage: true,
            confirmStage: false, 
            notification: `${course.class} is added to cart` 
        })
    }

    deleteCourse(course) {
        this.setState({
            enrolledCourses: this.state.enrolledCourses.filter(c => c._id !== course._id),
            notification: `${course.class} is removed from cart` 
        })
    }

    selectCoursesList() {
        return this.state.courses.map(course => {
            return <Course course={course} selectCourse={this.selectCourse} key={course._id}/>;
        })
    }

    enrolledCoursesList() {
        return this.state.enrolledCourses.map(course => {
            return <EnrolledCourse course={course}  deleteCourse={this.deleteCourse} key={course._id}/>;
        })
    }

    confirmedCourseList() {
        return this.state.enrolledCourses.map(course => {
            return <ConfirmedCourse course={course}  key={course._id}/>;
        })
    }

    SuccessCourseList() {
        return this.state.enrolledCourses.map(course => {
            return <SuccessCourse course={course} key={course._id}/>;
        })
    }

    // getCount(){
    //     fetch('http://localhost:3000/api/v1/student/1003123', {
    //         crossDomain:true,
    //         method: 'GET',
    //         headers: {'Content-Type':'application/json'},
    //         }).then(response => response.json()).then(data => this.setState({
    //             enrolledCoursesString: data.value
    //          }));
    // } 

    getCount(courseId){
        fetch('http://localhost:3001/read-from-node?courseid=' + courseId, {
            crossDomain:true,
            method: 'GET',
            headers: {'Content-Type':'text/plain'},
            })
            .then(function(response) {
                if (response.ok) {
                    // get count
                    this.setState({
                        count: response.text()
                    })
                }
                else {
                    // print error
                    console.log(response.text())
                }
            });
    }

     onSubmit(e) { 
        e.preventDefault()
        this.setState({
            selectStage: true, 
            searchStage: false,
            confirmStage: false, 
        })
    }

    render() {
        let step2Button
        let enrolmentsummary
        this.getCount()
        if (this.state.enrolledCourses.length != 0)  {
            step2Button =  <button type="button" className="btn btn-primary" onClick={this.goToStep2}>
                        Proceed to step 2 of 3
                    </button>
            enrolmentsummary = <div>
            <h3>Enrollment Summary</h3>
            <div></div>
            <table className="table">
            <thead className="thead-light">
                <tr>
                <th>Delete</th>
                <th>Class</th>
                <th>Days/Times</th>
                <th>Room</th>
                <th>Instructor</th>
                <th>Units</th>
                <th>Status</th>
                <th></th>
                </tr>
            </thead>
            <tbody>
                {this.enrolledCoursesList()}
            </tbody>
            </table>
            </div>
        } else {
            enrolmentsummary = <div className="summary">
                    <h2>No courses have been added to cart</h2>
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
                    <p> {notification} </p>
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
                            <input type="text" class="form-control" placeholder="Enter Course Number"  value={this.state.courseNumber}
                        onChange={this.onChangeCouseNumber}/>        
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
                        <div className="d-flex flex-row-reverse">
                            <div className="p-2">
                                {step2Button}
                            </div>
                        </div>
                </div>
            )
        }
        else if (this.state.selectStage) {
            return (
                <div>
                     <h3>Enrollment Summary</h3>
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
                         <th></th>
                         </tr>
                     </thead>
                     <tbody>
                         {this.selectCoursesList()}
                     </tbody>
                     </table>
                </div>
             )
        } 
        else if (this.state.confirmStage) {
            return (
                <div>
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
                            <th></th>
                            </tr>
                        </thead>
                        <tbody>
                            {this.confirmedCourseList()}
                        </tbody>
                        </table>
                        <div className="d-flex flex-row-reverse">
                            <div className="p-2">
                            <button type="button" className="btn btn-primary confirm-flex" onClick={this.cancelConfirm}>
                                Cancel
                            </button>
                            <button type="button" className="btn btn-primary confirm-flex" onClick={this.previousConfirm}>
                                    Previous 
                            </button>
                            <button type="button" className="btn btn-primary confirm-flex" onClick={this.confirmConfirm}>
                                    Finished Enrolling
                            </button>
                            </div>
                        </div>
                </div>
             )
        }
        else if (this.state.successStage) {
            return (
                <div>
                     <h3>View Results</h3>
                     <div></div>
                     <table className="table">
                     <thead className="thead-light">
                         <tr>
                         <th>Class</th>
                         <th>Message</th>
                         <th>Success</th>
                         </tr>
                     </thead>
                     <tbody>
                         {this.SuccessCourseList()}
                     </tbody>
                     </table>
                     <div className="d-flex flex-row-reverse">
                            <div className="p-2">
                                <button type="button" className="btn btn-primary confirm-flex" onClick={this.addAnotherClass}>
                                        Add Another Class 
                                </button>
                            </div>
                        </div>
                </div>
             )
        }
    }
}


