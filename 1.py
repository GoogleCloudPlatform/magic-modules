import os

def modify_yaml_files():
    products_dir = 'mmv1/products'
    
    # Fields to remove under each section
    remove_fields = {
        'operation': {'path', 'wait_ms', 'kind'},
        'result': {'path'},
        'status': {'path', 'allowed', 'complete'},
        'error': {'path', 'message'}
    }
    
    for product in os.listdir(products_dir):
        product_path = os.path.join(products_dir, product)
        if os.path.isdir(product_path):
            for root, _, files in os.walk(product_path):
                for file in files:
                    if file.endswith('.yaml'):
                        file_path = os.path.join(root, file)
                        print(f"\nProcessing: {file_path}")
                        
                        with open(file_path, 'r') as f:
                            lines = f.readlines()
                        
                        # Track state and indentation
                        in_async = False
                        current_section = None
                        async_indent = 0
                        section_indent = 0
                        lines_to_keep = []
                        
                        # Buffer for current section
                        current_section_lines = []
                        current_section_has_content = False
                        
                        for i, line in enumerate(lines, 1):
                            indent = len(line) - len(line.lstrip())
                            stripped = line.strip()
                            
                            # Check for async section
                            if stripped.startswith('async:'):
                                in_async = True
                                async_indent = indent
                                lines_to_keep.append(line)
                                continue
                            
                            # Check for sections under async
                            if in_async and any(stripped.startswith(f"{section}:") for section in remove_fields):
                                # Process previous section if exists
                                if current_section and current_section_has_content:
                                    lines_to_keep.extend(current_section_lines)
                                
                                current_section = next(section for section in remove_fields if stripped.startswith(f"{section}:"))
                                section_indent = indent
                                current_section_lines = [line]
                                current_section_has_content = False
                                continue
                            
                            # Process lines within a section
                            if current_section and indent > section_indent:
                                # Check if this is a field we want to remove
                                field = stripped.split(':')[0] if ':' in stripped else None
                                if field and field in remove_fields[current_section]:
                                    continue
                                
                                # If it's not a field to remove, keep track of content
                                if stripped and not stripped.startswith('#'):
                                    current_section_has_content = True
                                current_section_lines.append(line)
                            else:
                                # Leaving a section
                                if current_section:
                                    if current_section_has_content:
                                        lines_to_keep.extend(current_section_lines)
                                    current_section = None
                                    current_section_lines = []
                                    current_section_has_content = False
                                
                                # Check if we're leaving async section
                                if in_async and indent <= async_indent and stripped:
                                    in_async = False
                                
                                lines_to_keep.append(line)
                        
                        # Handle last section if exists
                        if current_section and current_section_has_content:
                            lines_to_keep.extend(current_section_lines)
                        
                        # Write back only if we made changes
                        if lines_to_keep != lines:
                            print(f"Writing changes to {file_path}")
                            with open(file_path, 'w') as f:
                                f.writelines(lines_to_keep)
                        else:
                            print(f"No changes needed in {file_path}")

if __name__ == "__main__":
    modify_yaml_files()