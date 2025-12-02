import { getContext, setContext } from 'svelte';
import * as VIAM from '@viamrobotics/sdk';

const ROBOT_KEY = Symbol('ROBOT');

enum SwitchLabels {
    off = "off",
    red = "red",
    green = "green",
    blue = "blue",
}

interface ReadingData {
    distance: number;
}

class RobotContext {
    private _sensor1Client: VIAM.SensorClient | undefined;
    private _switchClient: VIAM.SwitchClient | undefined;

    async initRobot() {
        const host = 'rgb-led-main.jqtpsxtyal.viam.cloud';

        const machine = await VIAM.createRobotClient({
            host,
            credentials: {
                type: 'api-key',

                payload: '9dluv6pvmohu31hl8mhrynilxj3t6urd',
                authEntity: '02977fdc-1cc1-41a2-b1a1-33f25250e5cf',

            },
            signalingAddress: 'https://app.viam.com:443',
        });

        // sensor-1
        this._sensor1Client = new VIAM.SensorClient(machine, 'motion-sensor');

        // switch-1
        this._switchClient = new VIAM.SwitchClient(machine, 'rgb');
    }

    public async getDistanceReading(): Promise<ReadingData> {
        const sensor1ReturnValue = await this._sensor1Client?.getReadings() || { distance: 0 };
        return sensor1ReturnValue as unknown as ReadingData;
    }

    public async getSwitchState(): Promise<SwitchLabels> {
        const switchReturnValue = await this._switchClient?.getPosition() || 0;
        const [, switchLabels] = await this._switchClient?.getNumberOfPositions() || [0, ["off","red", "green", "blue"]];

        return switchLabels[switchReturnValue] as SwitchLabels;
    }

    public async setSwitchState(state: number): Promise<void> {
        await this._switchClient?.setPosition(state);
    }
}

export function createRobotContext() {
    const robotContext = new RobotContext();
    robotContext.initRobot();

    setContext(ROBOT_KEY, robotContext);
    return robotContext;
}

export function useRobot() {
    return getContext<RobotContext>(ROBOT_KEY);
}
