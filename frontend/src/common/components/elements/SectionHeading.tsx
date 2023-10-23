import React, { ReactNode } from "react";

interface SectionHeadingProps {
  title: string,
  icon: ReactNode,
  className?: string,
}

export const SectionHeading = ({ title, icon, className = '' }: SectionHeadingProps) => {
  return (
    <div className={`flex items-center gap-4 ${className}`}>
      {icon && <>{icon}</>}
      <h2 className="capitalize">{title}</h2>
    </div>
  )
}