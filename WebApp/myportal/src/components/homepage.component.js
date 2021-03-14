import React, { Component } from 'react';
import Course from "./course.component";

export default class HomePage extends Component {
    constructor(props) {
        super(props); 

        this.onChangeSubject = this.onChangeSubject.bind(this);
        this.onChangeCouseNumber = this.onChangeCouseNumber.bind(this);
        this.onChangeCourseCareer = this.onChangeCourseCareer.bind(this);
        this.onSubmit = this.onSubmit.bind(this);

        this.state = {
            subjects : ['01', '02', '03', '10', '20', '30', '40', '50', '51', '99'], 
            courseNums: ['contains', 'greater than or equal to', 'is exactly', 'lesser than or equal to'],
            courseCareers: ['Continuing Education Training', 'Master', 'Master of Architecture', 'Non-Graduating', 'PhD', "SUTD MIT Master", 'Undergradutate'],
            subject: '', 
            number: '', 
            career: '', 
            courses: []
            // search: false,
        }
        
    }

    componentDidMount() {
        const course = this.props.location.query
        const cart = []
        if (course != undefined) {
            cart.push(course.course)
            this.setState({
                courses: cart
            })
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

     coursesList() {
        return this.state.courses.map(course => {
            // console.log(course)
            return <Course course={course} key={course._id}/>;
        })
    }

     onSubmit(e) {
         // check for course
        this.props.history.push("/select" );   
    }

    render() {
        const cart = this.state.courses
        let addedToCart
        if (cart != "" ) {
            addedToCart = <p>{this.state.courses[0].class} has been added to your cart</p> 
        } else {
            addedToCart = <div></div>
        }
        return (
            <div>
                {addedToCart}
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
                        {/* <th>Delete</th> */}
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
                        {this.coursesList()}
                    </tbody>
                    </table>
                    </div>
                </div>
                </div> 
            </div>
        )
    }
}

