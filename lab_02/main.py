import sys
import copy
import numpy as np
from scipy.integrate import odeint
from PyQt5 import QtWidgets, QtCore
from matplotlib import pyplot as plt


def calculate_steady_state(intensity_matrix):
    num_states = intensity_matrix.shape[0]
    q_matrix = copy.deepcopy(intensity_matrix.T)
    q_matrix[-1, :] = 1
    b = np.zeros(num_states)
    b[-1] = 1
    steady_state = np.linalg.solve(q_matrix, b)
    return steady_state


def calculate_stabilization_time(intensity_matrix, steady_state):
    intensity_matrix = np.transpose(intensity_matrix)

    def system(y, t, matrix):
        return matrix @ y

    initial_state = np.zeros(len(steady_state))
    initial_state[0] = 1

    time = np.linspace(0, 100, 1000)

    solution = np.transpose(odeint(system, initial_state, time, args=(intensity_matrix,)))

    stabilization_times = [-1 for _ in range(len(steady_state))]

    for i in range(len(steady_state)):
        plt.plot(time, solution[i])
    plt.show()

    for i in range(len(steady_state)):
        for j in range(len(solution[i]) - 1, -1, -1):
            if np.abs(solution[i][j] - steady_state[i]) > 0.01:
                stabilization_times[i] = time[j]
                break

    return stabilization_times


class SteadyStateCalculator(QtWidgets.QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("Лабораторная работа №2")
        self.setStyleSheet("background-color: #f0f0f0;")

        # Layouts
        layout = QtWidgets.QVBoxLayout()
        grid_layout = QtWidgets.QGridLayout()

        # Поля для ввода количества состояний
        self.state_count_label = QtWidgets.QLabel("Количество состояний (≤10):")
        self.state_count_label.setStyleSheet("color: #333333; font-size: 12px;")
        layout.addWidget(self.state_count_label)

        self.state_count_entry = QtWidgets.QLineEdit()
        self.state_count_entry.setFixedWidth(100)
        self.state_count_entry.setStyleSheet("font-size: 12px;")
        self.state_count_entry.textChanged.connect(self.update_matrix_entries)
        layout.addWidget(self.state_count_entry)

        # Метка для матрицы интенсивностей
        self.matrix_label = QtWidgets.QLabel("Матрица интенсивностей переходов:")
        self.matrix_label.setStyleSheet("color: #333333; font-size: 12px;")
        layout.addWidget(self.matrix_label)

        # Поле для ввода матрицы интенсивностей переходов
        self.matrix_entries = []
        for i in range(10):
            row_entries = []
            for j in range(10):
                entry = QtWidgets.QLineEdit("0.0")
                entry.setFixedWidth(50)
                entry.setAlignment(QtCore.Qt.AlignCenter)
                entry.setStyleSheet("font-size: 12px; border: 1px solid #333;")
                entry.setEnabled(False)
                entry.installEventFilter(self)  # Установка фильтра событий для каждой ячейки
                grid_layout.addWidget(entry, i, j)
                row_entries.append(entry)
            self.matrix_entries.append(row_entries)

        layout.addLayout(grid_layout)

        # Кнопка для расчета решения
        self.solve_button = QtWidgets.QPushButton("Вычислить")
        self.solve_button.setStyleSheet("font-size: 12px; background-color: #4CAF50; color: white;")
        self.solve_button.clicked.connect(self.solve)
        layout.addWidget(self.solve_button)

        # Поле для вывода результата
        self.result_text = QtWidgets.QTextEdit()
        self.result_text.setFixedHeight(200)
        self.result_text.setStyleSheet("font-size: 12px; border: 1px solid #333;")
        layout.addWidget(self.result_text)

        self.setLayout(layout)

        # Установка клавиши Enter для кнопки "Вычислить"
        self.solve_button.setDefault(True)
        self.solve_button.setAutoDefault(True)

    def update_matrix_entries(self):
        try:
            num_states = int(self.state_count_entry.text())
            if num_states > 10:
                raise ValueError("Количество состояний не должно превышать 10.")
            for i in range(10):
                for j in range(10):
                    entry = self.matrix_entries[i][j]
                    if i < num_states and j < num_states:
                        entry.setEnabled(True)
                    else:
                        entry.setText("0.0")
                        entry.setEnabled(False)
        except ValueError:
            pass

    def solve(self):
        try:
            num_states = int(self.state_count_entry.text())
        except ValueError:
            QtWidgets.QMessageBox.critical(self, "Ошибка", "Количество состояний должно быть целым числом ≤10.")
            return

        if num_states > 10:
            QtWidgets.QMessageBox.critical(self, "Ошибка", "Количество состояний должно быть целым числом ≤10.")
            return

        intensity_matrix = np.zeros((num_states, num_states))
        for i in range(num_states):
            for j in range(num_states):
                intensity_matrix[i, j] = float(self.matrix_entries[i][j].text().replace(",", "."))

        for i in range(num_states):
            intensity_matrix[i, i] = -np.sum(intensity_matrix[i]) + intensity_matrix[i, i]

        try:
            steady_state_probabilities = calculate_steady_state(intensity_matrix)
            stabilization_times = calculate_stabilization_time(intensity_matrix, steady_state_probabilities)
        except np.linalg.LinAlgError:
            QtWidgets.QMessageBox.critical(self, "Ошибка", "Заданы некорректные значения.")
            return

        result = "Стабилизированные вероятности для каждого состояния:\n"
        result += "\n".join([f"Состояние {i}: {steady_state_probabilities[i]:.4f}" for i in range(num_states)])
        result += "\n\nВремя стабилизации для каждого состояния:\n"
        result += "\n".join([f"Состояние {i}: {stabilization_times[i]:.4f}" for i in range(num_states)])

        self.result_text.setPlainText(result)

    def eventFilter(self, obj, event):
        if event.type() == QtCore.QEvent.KeyPress:
            key = event.key()
            for i in range(10):
                for j in range(10):
                    if obj == self.matrix_entries[i][j]:
                        if key == QtCore.Qt.Key_Right and j < 9:
                            self.matrix_entries[i][j + 1].setFocus()
                            return True
                        elif key == QtCore.Qt.Key_Left and j > 0:
                            self.matrix_entries[i][j - 1].setFocus()
                            return True
                        elif key == QtCore.Qt.Key_Down and i < 9:
                            self.matrix_entries[i + 1][j].setFocus()
                            return True
                        elif key == QtCore.Qt.Key_Up and i > 0:
                            self.matrix_entries[i - 1][j].setFocus()
                            return True
        return super().eventFilter(obj, event)


if __name__ == "__main__":
    app = QtWidgets.QApplication(sys.argv)
    window = SteadyStateCalculator()
    window.show()
    sys.exit(app.exec_())
