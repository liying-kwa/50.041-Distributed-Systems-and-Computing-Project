import React from 'react';

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

export default Course
