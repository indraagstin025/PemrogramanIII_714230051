import { Button } from "@material-tailwind/react";

export function ButtonAtom({ children, onClick, color = "blue", className = "", type = "button"}) {
  return (
    <Button color={color} className={className} onClick={onClick} type={type}>
      {children}
    </Button>
  );
}