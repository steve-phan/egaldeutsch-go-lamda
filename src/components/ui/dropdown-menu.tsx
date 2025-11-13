import React, { useState, useRef, useEffect } from "react";
import { cn } from "@/lib/utils";
import { ChevronDown } from "lucide-react";

interface DropdownMenuProps {
  trigger: React.ReactNode;
  children: React.ReactNode;
  align?: "left" | "right";
  className?: string;
}

export const DropdownMenu: React.FC<DropdownMenuProps> = ({
  trigger,
  children,
  align = "left",
  className,
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div className={cn("relative", className)} ref={dropdownRef}>
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-1 text-foreground hover:text-primary px-3 py-2 rounded-md text-sm font-medium transition-colors"
      >
        {trigger}
        <ChevronDown
          size={16}
          className={cn(
            "transition-transform",
            isOpen && "transform rotate-180"
          )}
        />
      </button>

      {isOpen && (
        <div
          className={cn(
            "absolute top-full mt-2 w-48 bg-card rounded-md shadow-lg border border-border z-50",
            align === "right" ? "right-0" : "left-0"
          )}
        >
          <div
            className="py-1"
            onClick={() => setIsOpen(false)}
          >
            {children}
          </div>
        </div>
      )}
    </div>
  );
};

interface DropdownMenuItemProps {
  children: React.ReactNode;
  onClick?: () => void;
  className?: string;
}

export const DropdownMenuItem: React.FC<DropdownMenuItemProps> = ({
  children,
  onClick,
  className,
}) => {
  return (
    <div
      onClick={onClick}
      className={cn(
        "px-4 py-2 text-sm text-foreground hover:bg-accent hover:text-accent-foreground cursor-pointer transition-colors",
        className
      )}
    >
      {children}
    </div>
  );
};
