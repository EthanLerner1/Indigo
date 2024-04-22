import tkinter as tk
from tkinter import messagebox, scrolledtext
from tkinter.ttk import Combobox
import socket
import threading

class CommanderApp:
    def __init__(self, master):
        self.master = master
        master.title("Commander")

        self.label = tk.Label(master, text="Select destination client:")
        self.label.grid(row=0, column=0, padx=5, pady=5, sticky="w")

        self.client_combo = Combobox(master, width=30, state="readonly")
        self.client_combo.grid(row=0, column=1, padx=5, pady=5, sticky="ew")

        self.refresh_button = tk.Button(master, text="Refresh", command=self.refresh_clients)
        self.refresh_button.grid(row=0, column=2, padx=5, pady=5)

        self.command_text = scrolledtext.ScrolledText(master, wrap=tk.WORD, width=60, height=10)
        self.command_text.grid(row=1, column=0, columnspan=3, padx=5, pady=5, sticky="ew")

        self.command_text = scrolledtext.ScrolledText(master, wrap=tk.WORD, width=60, height=10)
        self.command_text.grid(row=1, column=0, columnspan=3, padx=5, pady=5, sticky="ew")

        self.result_text = scrolledtext.ScrolledText(master, wrap=tk.WORD, width=60, height=10)
        self.result_text.grid(row=2, column=0, columnspan=3, padx=5, pady=5, sticky="ew")

        self.send_button = tk.Button(master, text="Send Command", command=self.send_command)
        self.send_button.grid(row=3, column=0, columnspan=3, padx=5, pady=5, sticky="ew")

        self.refresh_clients()

    def refresh_clients(self):
        try:
            with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
                s.connect(("localhost", 9999))
                s.sendall(b"LIST\n")
                response = s.recv(4096).decode()
                clients = response.splitlines()
                self.client_combo["values"] = clients
        except Exception as e:
            messagebox.showerror("Error", f"Failed to refresh clients: {e}")

    def send_command(self):
        destination_client = self.client_combo.get()
        command = self.command_text.get("1.0", tk.END).strip()
        try:
            with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
                s.connect(("localhost", 9999))
                s.sendall(f"{destination_client}|{command}\n".encode())
                response = s.recv(4096).decode()
                print(response)
                self.result_text.insert(tk.END, f"> {command}\n{response}\n")
                self.result_text.see(tk.END)  # Scroll to the bottom
        except Exception as e:
            messagebox.showerror("Error", f"Failed to send command: {e}")

def main():
    root = tk.Tk()
    app = CommanderApp(root)
    root.mainloop()

if __name__ == "__main__":
    main()
