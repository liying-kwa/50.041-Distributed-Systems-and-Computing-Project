import React, { Component } from 'react';

const Course  = (props) => {
    return (
    <tr>
      <td>{props.course.class}</td>
      <td>{props.course.dayTime}</td>
      <td>{props.course.room}</td>
      <td>{props.course.instructor}</td>
      <td>{props.course.units}</td>
      <td>{props.course.status}</td>
    </tr>
    )
}

export default Course; 
