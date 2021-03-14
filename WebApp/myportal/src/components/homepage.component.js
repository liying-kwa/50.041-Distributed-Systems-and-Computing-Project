import React, { Component } from 'react';

const Course = props => (
    <tr>
      <td>{props.course.class}</td>
      <td>{props.course.section}</td>
      <td>{props.course.dayTime}</td>
      <td>{props.course.room}</td>
      <td>{props.course.instructor}</td>
      <td>{props.course.meetDate}</td>
      <td>{props.course.status}</td>
      <td>
        <button type="button" className="btn btn-primary" onClick={() => {props.selectCourse(props.course)}}>
                Select
        </button>
      </td>
    </tr>
)

const EnrolledCourse = props => (
    <tr>
    <td>
        <button type="button" className="btn btn-primary" onClick={() => {props.deleteCourse(props.course)}}>
            Delete
        </button>
      </td>
      <td>{props.course.section} ({props.course.class})</td>
      <td>{props.course.dayTime}</td>
      <td>{props.course.room}</td>
      <td>{props.course.instructor}</td>
      <td>{props.course.units}</td>
      <td>{props.course.status}</td>
    </tr>
)


export default class HomePage extends Component {
    constructor(props) {
        super(props); 

        this.onChangeSubject = this.onChangeSubject.bind(this);
        this.onChangeCouseNumber = this.onChangeCouseNumber.bind(this);
        this.onChangeCourseCareer = this.onChangeCourseCareer.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
        this.selectCourse = this.selectCourse.bind(this);
        this.deleteCourse = this.deleteCourse.bind(this);

        this.state = {
            subjects : ['01', '02', '03', '10', '20', '30', '40', '50', '51', '99'], 
            courseNums: ['contains', 'greater than or equal to', 'is exactly', 'lesser than or equal to'],
            courseCareers: ['Continuing Education Training', 'Master', 'Master of Architecture', 'Non-Graduating', 'PhD', "SUTD MIT Master", 'Undergradutate'],
            subject: '', 
            number: '', 
            career: '', 
            enrolledCourses: [],
            searchStage: true,
            selectStage: false, 
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

    selectCourse(course) {
        const enrolledCourses = [...this.state.enrolledCourses];
        enrolledCourses.push(course)
        this.setState({
            enrolledCourses: enrolledCourses,
            selectStage: false, 
            searchStage: true,
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

     onSubmit(e) { 
        e.preventDefault()
        this.setState({
            selectStage: true, 
            searchStage: false
        })
    }

    render() {
        const notification = this.state.notification
        if (this.state.searchStage) {
            return (
                <div>
                    <p> {notification} </p>
                    <div className = 'row'>
                    <div className="col-4">
                    <h2>Class Search</h2> 
                    <form onSubmit={this.onSubmit}>
                        <div className = "form-group"> 
                            <label>Subject</label>
                            <select ref="userInput"
                                required
                                className="form-control"
                                value={this.state.subject}
                                onChange={this.onChangeSubject}>
                                {
                                    this.state.subjects.map(function(subject) {
                                        return <option
                                        key={subject}
                                        value={subject}>{subject}
                                        </option>;
                                    })
                                }
                            </select>
                        </div>
                        <div className = "form-group"> 
                            <label>Course Number</label>
                            <select ref="userInput"
                                required
                                className="form-control"
                                value={this.state.subject}
                                onChange={this.onChangeSubject}>
                                {
                                    this.state.courseNums.map(function(number) {
                                        return <option
                                        key={number}
                                        value={number}>{number}
                                        </option>;
                                    })
                                }
                            </select>
                        </div>
                        <div className = "form-group"> 
                            <label>Course Career</label>
                            <select ref="userInput"
                                required
                                className="form-control"
                                value={this.state.subject}
                                onChange={this.onChangeSubject}>
                                {
                                    this.state.courseCareers.map(function(career) {
                                        return <option
                                        key={career}
                                        value={career}>{career}
                                        </option>;
                                    })
                                }
                            </select>
                        </div>
                            <div className="form-group">
                                <input type="submit" value="Search" className="btn btn-primary"/>
                            </div>
                    </form>
                    </div>
                    <div className="col-8">
                    <div>
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
    }
}


