import React, { Component } from 'react';
import { Link } from 'react-router-dom';

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
        <Link to={{pathname: `/`, query: {class: props.course.class}}}>
            <button type="button" className="btn btn-primary">
                Select
            </button>
        </Link> 
      </td>
    </tr>
)



export default class SelectPage extends Component {
    constructor(props) {
        super(props); 
            this.state = {
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
                    }
                ],
            }
        }

    coursesList() {
        return this.state.courses.map(course => {
            return <Course course={course} key={course._id}/>;
        })
    }
    
    render() {
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
                    {this.coursesList()}
                </tbody>
                </table>
           </div>
        )
    }
}

