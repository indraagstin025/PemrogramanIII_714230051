import { Input } from "@material-tailwind/react";

export function TextFieldAtom({ size, color, label, name, value, onChange, error = false, helperText = '', type = 'text' }) {
    return (
        <div className="w-full">
            <Input
                size={size || 'lg'}
                color={color || 'blue'}
                label={label}
                name={name}
                value={value}
                onChange={onChange}
                type={type}
                error={error}
                className="w-full"
            />
            {error && <p className="text-red-500 text-xs mt-1">{helperText}</p>}
        </div>
    );
}