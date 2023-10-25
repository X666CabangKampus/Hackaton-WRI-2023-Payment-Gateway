import React from "react"

interface SubHeadingProps {
  children: String,
  className?: React.ReactNode,
}

export const SectionSubHeading = ({ children, className }: SubHeadingProps) => {
  return (
    <div className={`${className}`}>
      <p>{children}</p>
    </div>
  )
}