import React from 'react';

const EnrolledCourse = props => (
    <tr>
      <td>{props.course.section} ({props.course.class})</td>
      <td>{props.course.dayTime}</td>
      <td>{props.course.room}</td>
      <td>{props.course.instructor}</td>
      <td>{props.course.units}</td>
      <td>{props.course.status}</td>
    </tr>
)

export default EnrolledCourse;