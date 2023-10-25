import { ReactNode } from "react"

interface ContainerProps {
  children?: ReactNode
}


export const Container = ({ children }: ContainerProps) => {
  return (
    <div>
      {children}
    </div>
  )
}